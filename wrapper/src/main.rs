use wasmtime::component::{*, self};
use wasmtime::{Config, Engine, Store};

bindgen!(in "../producer/wit");


fn main() -> wasmtime::Result<()> {
    let mut config = Config::new();
    config.wasm_component_model(true);
    let engine = Engine::new(&config);
    //let component = Component::

    Ok(())
}
