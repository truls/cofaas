use wasmtime::component::*;
use wasmtime::{Config, Engine, Store};
use wasmtime_wasi::preview2;
use wasmtime_wasi::preview2::{Table, WasiCtx, WasiCtxBuilder, WasiView};

use crate::exports::chained_service::api::hello_world::HelloRequest;

bindgen!({
    path:"../wit",
    async: true,
    world: "top-level",
});

struct MyState {
    table: Table,
    wasi: WasiCtx,
}

impl WasiView for MyState {
    fn table(&self) -> &Table {
        &self.table
    }
    fn table_mut(&mut self) -> &mut Table {
        &mut self.table
    }
    fn ctx(&self) -> &WasiCtx {
        &self.wasi
    }
    fn ctx_mut(&mut self) -> &mut WasiCtx {
        &mut self.wasi
    }
}

#[tokio::main]
async fn main() -> wasmtime::Result<()> {
    let mut config = Config::new();
    config.wasm_component_model(true);
    config.async_support(true);
    let engine = Engine::new(&config)?;
    //let component = Component::from_file(&engine, "../producer/component/main-component.wasm")?;
    let component = Component::from_file(&engine, "../producer/component/composed.wasm")?;

    let mut linker = Linker::new(&engine);

    let mut table = Table::new();
    let wasi = WasiCtxBuilder::new().inherit_stdio().build(&mut table)?;

    preview2::wasi::command::add_to_linker(&mut linker)?;

    let mut store = Store::new(&engine, MyState { table, wasi });
    let (bindings, _) =
        TopLevel::instantiate_async(&mut store, &component, &linker).await?;

    let hello_world = bindings.chained_service_api_hello_world();
    hello_world.call_init_component(&mut store).await?;
    let arg = HelloRequest { name: &"foo" };
    let res = hello_world.call_say_hello(&mut store, arg).await?;

    println!("{:?}", res);

    Ok(())
}
