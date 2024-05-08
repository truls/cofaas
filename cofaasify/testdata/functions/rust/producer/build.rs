fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("helloworld.proto")?;
    tonic_build::compile_protos("prodcon.proto")?;
    Ok(())
}
