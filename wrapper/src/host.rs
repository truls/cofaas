use std::path::PathBuf;

use exports::cofaas::application::greeter::{HelloReply, HelloRequest};
use wasmtime::component::*;
use wasmtime::component::__internal::anyhow::anyhow;
use wasmtime::{Config, Engine, Store};
use wasmtime_wasi::preview2;
use wasmtime_wasi::preview2::{Table, WasiCtx, WasiCtxBuilder, WasiView};


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

pub struct MyComponent {
    bindings: Box<(TopLevel, Instance)>,
    store: Box<Store<MyState>>,
}

fn get_env_or_default(env: &str) -> String {
    match std::env::var(env) {
        Ok(val) => val,
        Err(_) => "1".to_string(),
    }
}

impl MyComponent {
    pub async fn new(component: PathBuf) -> wasmtime::Result<MyComponent> {
        let mut config = Config::new();
        config.wasm_component_model(true);
        config.async_support(true);
        let engine = Engine::new(&config)?;

        let component = Component::from_file(&engine, component)?;

        let mut linker = Linker::new(&engine);

        let mut table = Table::new();
        let wasi = WasiCtxBuilder::new()
            .inherit_stdio()
            .push_env("TRANSFER_SIZE_KB", get_env_or_default("TRANSFER_SIZE_KB"))
            .push_env("REPEATS", get_env_or_default("REPEATS"))
            .build(&mut table)?;

        preview2::wasi::command::add_to_linker(&mut linker)?;

        let mut store = Box::new(Store::new(&engine, MyState { table, wasi }));

        let instance =
            Box::new(TopLevel::instantiate_async(&mut *store, &component, &linker).await?);

        let mut this = Self{bindings: instance,
                        store};
        let iface = this.bindings.0.cofaas_application_greeter();

        iface.call_init_component(&mut *this.store).await?;

        Ok(this)
    }

    pub async fn call(&mut self) -> wasmtime::Result<HelloReply> {
        let arg = HelloRequest { name: "foo".to_string() };
        let iface = self.bindings.0.cofaas_application_greeter();
        let res = iface.call_say_hello(&mut *self.store, &arg).await?;
        res.or_else(|x| Err(anyhow!("Call failed with exit code {}", x)))
    }
}
