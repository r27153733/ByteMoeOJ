use crate::pb::wasm::wasm_executor_server::WasmExecutor;
use crate::pb::wasm::{WasmExecutionOutput, WasmExecutionRequest, WasmExecutionResponse};
use async_trait::async_trait;
use tonic::{Request, Response, Status};
use crate::wasm_run::run::{err_code, run_wasm};

// Implementation of the WasmExecutor trait.
pub struct WasmServer {
    // You can add necessary state here if required.
}

impl WasmServer {
    pub fn new() -> Self {
        WasmServer {}
    }
}

#[async_trait]
impl WasmExecutor for WasmServer {
    async fn execute(
        &self,
        request: Request<WasmExecutionRequest>,
    ) -> Result<Response<WasmExecutionResponse>, Status> {
        let req = request.into_inner();
        let mut resp = WasmExecutionResponse {
            outputs: vec![],
        };
        for wasm_binary in req.wasm_binary_arr {
            for input in &req.inputs {
                let wasm_binary = wasm_binary.as_slice();
                let stdin = input.stdin.as_slice();
                let args = input.args.as_slice();
                let vec: Vec<(String, String)>;
                if input.envs.is_empty() {
                    vec = vec![];
                } else {
                    vec = input.envs.iter()
                        .map(|(k, v)| (k.clone(), v.clone()))
                        .collect();
                }
                
                let envs: &[(String, String)] = &vec;
                let memory_limit = input.memory_limit;
                let fuel_limit = input.fuel_limit;
                let stdout_limit = input.stdout_limit;
                let stderr_limit = input.stderr_limit;

                let res = run_wasm(wasm_binary, stdin, args, envs,
                                   memory_limit as usize, fuel_limit, stdout_limit, stderr_limit);
                match res {
                    Ok(res) => {
                        let hash = xxhash_rust::const_xxh64::xxh64(res.stdout.as_slice(), 0);

                        let input_data = res.stdout.as_slice();
                        let cleaned_data: Vec<u8> = input_data.iter()
                            .filter(|&&x| x != b' ' && x != b'\n' && x != b'\r') // 过滤掉空格、换行符和回车符
                            .cloned()
                            .collect();
                        let token_hash = xxhash_rust::const_xxh64::xxh64(cleaned_data.as_slice(), 0);


                        resp.outputs.push(WasmExecutionOutput{
                            stdout: res.stdout,
                            stderr: res.stderr,
                            memory_used: res.memory_used as u64,
                            fuel_consumed: res.fuel_consumed,
                            status: 5,
                            stdout_hash: hash,
                            stdout_token_stream_hash: token_hash,
                        })
                    }
                    Err(err) => {
                        let str = err.to_string();
                        
                        resp.outputs.push(WasmExecutionOutput{
                            stdout: vec![],
                            stderr: str.into(),
                            memory_used: 0,
                            fuel_consumed: 0,
                            status: err_code(err),
                            stdout_hash: 0,
                            stdout_token_stream_hash: 0,
                        })
                    }
                }
            }
        }
        
        Ok(Response::new(resp))
    }
}

