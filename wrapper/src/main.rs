use wasmtime::component::*;
use wasmtime::{Config, Engine, Store};

bindgen!(in "../producer/wit");

fn main() -> wasmtime::Result<()> {
    println!("Hello, world!");

    Ok(())
}
