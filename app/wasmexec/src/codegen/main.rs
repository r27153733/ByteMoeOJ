fn main() {
    tonic_build::configure()
        .build_client(false)
        .out_dir("src/pb")
        .compile_protos(&["src/pb/wasmexecutor.proto"], &["src/pb"])
        .expect("failed to compile protos");
}