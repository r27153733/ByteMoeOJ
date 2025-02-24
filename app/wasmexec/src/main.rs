use tonic::codec::CompressionEncoding;

use crate::pb::wasm::wasm_executor_server::WasmExecutorServer;
use crate::server::wasm_server::WasmServer;

mod server;
mod pb;
mod wasm_run;

// Create a gRPC server to serve the WasmExecutor service.
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let wasm_server = WasmServer::new();
    let svc = WasmExecutorServer::new(wasm_server)
    .max_encoding_message_size(usize::MAX)
    .max_decoding_message_size(usize::MAX)
    .accept_compressed(CompressionEncoding::Zstd)
    .send_compressed(CompressionEncoding::Zstd);

    println!("Server listening on {}", addr);

    tonic::transport::Server::builder()
        .add_service(svc)
        .serve(addr)
        .await?;

    Ok(())
}
