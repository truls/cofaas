fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("../cofaasify/testdata/protos/helloworld.proto")?;
    Ok(())
}
