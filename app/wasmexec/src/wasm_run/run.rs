use anyhow::Error;
use lazy_static::lazy_static;
use wasi_common::pipe::WritePipe;
use wasi_common::sync::WasiCtxBuilder;
use wasi_common::WasiCtx;
use wasmtime::{Config, Engine, Linker, Module, Store, Trap};
use crate::wasm_run::limit_err::LimitErr;
use crate::wasm_run::limited_buffer::LimitedBuffer;
use crate::wasm_run::memory_limiter::MemoryLimiter;

lazy_static! {
    static ref WasmEngine: Engine = {
        let mut config = Config::default();
        config.consume_fuel(true);
        Engine::new(&config).unwrap()
    };
}

struct State {
    limits: MemoryLimiter,
    wasi: WasiCtx,
}

pub struct WasmRunResult {
    pub stdout: Vec<u8>,
    pub stderr: Vec<u8>,
    pub memory_used: usize,
    pub fuel_consumed: u64,
}

pub fn run_wasm(
    wasm_binary: &[u8],
    stdin_bytes: &[u8],
    args: &[String],
    envs: &[(String, String)],
    memory_limit: usize,
    fuel_limit: u64,
    stdout_limit: u64,
    stderr_limit: u64,
) -> Result<(WasmRunResult), anyhow::Error> {
    let stdout = Box::new(WritePipe::new(LimitedBuffer::new(stdout_limit as usize)));
    let stderr = Box::new(WritePipe::new(LimitedBuffer::new(stderr_limit as usize)));
    let stdin = Box::new(wasi_common::pipe::ReadPipe::from(stdin_bytes.to_vec()));

    let memory_used;
    let remaining_fuel;
    let res: Result<(), Error>;
    {
        // let mut config = Config::default();
        // config.consume_fuel(true);
        // let engine = Engine::new(&config)?;
        let engine = &WasmEngine;
        let resource_limiter = MemoryLimiter::new(memory_limit);
        let mut linker = Linker::new(&engine);
        wasi_common::sync::add_to_linker(&mut linker, |state: &mut State| &mut state.wasi)?;

        let mut wasi_builder = WasiCtxBuilder::new();
        wasi_builder.stdout(stdout.clone()).stderr(stderr.clone());
        
        if !stdin_bytes.is_empty() {
            wasi_builder.stdin(stdin);
        }
        if !args.is_empty() {
            wasi_builder.args(args)?;
        }
        if !envs.is_empty() {
           wasi_builder.envs(envs)?;
        }
        let wasi = wasi_builder.build();
        
        let my_state = State {
            limits: resource_limiter,
            wasi,
        };

        let mut store = Store::new(&engine, my_state);
        store.limiter(|state| &mut state.limits);
        store.set_fuel(fuel_limit)?;

        let module = Module::new(&engine, wasm_binary)?;
        linker.module(&mut store, "", &module)?;
        
        res = linker
            .get_default(&mut store, "")?
            .typed::<(), ()>(&store)?
            .call(&mut store, ());
        
        memory_used = store.data().limits.get_try_alloc_size();
        remaining_fuel = store.get_fuel()?;
    }

        // 这里不可能有其他引用，所以不会失败。
        let stdout_buf = stdout.try_into_inner().unwrap();
        let stderr_buf = stderr.try_into_inner().unwrap();

    match res {
        Err(error) => {
            if memory_used > memory_limit {
                return Err(LimitErr::MemoryLimitExceeded.into());
            } else if stdout_buf.get_try_write_size() > stdout_limit as usize {
                return Err(LimitErr::StdoutOutputLimitExceeded.into());
            } else if stderr_buf.get_try_write_size() > stderr_limit as usize{
                return Err(LimitErr::StderrOutputLimitExceeded.into());
            } else if let Some(exit_code) = error.downcast_ref::<wasi_common::I32Exit>() {
                // exit code 为 0 是正常的。
                if exit_code.0 != 0 {
                    return Err(anyhow::anyhow!("exit code {}", exit_code.0));
                }
            } else if let Some(trap) = error.downcast_ref::<Trap>() {
                return match trap {
                    Trap::OutOfFuel => Err(LimitErr::FuelLimitExceeded.into()),
                    // Trap::MemoryOutOfBounds => Err(error),
                    _ => Err(error),
                }
            } else {
                return Err(error);
            }
        }
        _ => {}
    }

    Ok(WasmRunResult{
        stdout: stdout_buf.into_buffer(),
        stderr: stderr_buf.into_buffer(),
        memory_used,
        fuel_consumed: fuel_limit - remaining_fuel,
    })
}

pub fn err_code(err: anyhow::Error) -> u32 {
    if let Some(limit_err) = err.downcast_ref::<LimitErr>() {
        match limit_err {
            LimitErr::MemoryLimitExceeded => 7,
            LimitErr::FuelLimitExceeded => 8,
            LimitErr::StdoutOutputLimitExceeded => 9,
            LimitErr::StderrOutputLimitExceeded => 9,
        }
    } else {
        // runtime err
        10
    }
}