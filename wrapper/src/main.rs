use wasmtime::component::*;
use wasmtime::{Config, Engine, Store};
use wasmtime_wasi::preview2::{WasiCtx, WasiCtxBuilder, Table, WasiView};
use wasmtime_wasi::preview2;

use crate::exports::hello_world::HelloRequest;

bindgen!({
    path:"../producer/component/wit",
    async: true
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
    let component = Component::from_file(&engine, "../producer/component/main-component.wasm")?;

    let mut linker = Linker::new(&engine);

    let mut table = Table::new();
    //wasmtime_wasi::
    let wasi = WasiCtxBuilder::new().inherit_stdio().set_args(&[""]).build(&mut table)?;

    // preview2::wasi::filesystem::filesystem::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::io::streams::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::environment::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::preopens::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::exit::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::stdin::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::stdout::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::cli_base::stderr::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::clocks::wall_clock::add_to_linker(&mut linker, |x| x)?;
    // preview2::wasi::clocks::monotonic_clock::add_to_linker(&mut linker, |x| x)?;

    preview2::wasi::command::add_to_linker(&mut linker)?;

    let mut store = Store::new(
        &engine,
        MyState {  table, wasi },
    );
    let (bindings, _) = ProducerInterface::instantiate_async(&mut store, &component, &linker).await?;

    let interf = bindings.hello_world();
    interf.call_test(&mut store).await?;
    interf.call_init_component(&mut store).await?;

    let arg = HelloRequest{name: &"foo"};
    let res = interf.call_say_hello(&mut store, arg).await?;

    println!("{:?}", res);

    Ok(())
}
