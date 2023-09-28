mod host;
use host::MyComponent;
use tonic::{transport::Server, Request, Response, Status};

use hello_world::greeter_server::{Greeter, GreeterServer};
use hello_world::{HelloReply, HelloRequest};

use futures::lock::Mutex;
use std::net::SocketAddr;
use std::path::PathBuf;
use std::sync::Arc;
use clap::Parser;
use std::net::IpAddr;
use std::net::Ipv6Addr;

pub mod hello_world {
    tonic::include_proto!("helloworld"); // The string specified here must match the proto package name
}

pub struct MyGreeter {
    c: Arc<Mutex<MyComponent>>,
}

const LOCALHOST_V6: IpAddr = IpAddr::V6(Ipv6Addr::new(0, 0, 0, 0, 0, 0, 0, 1));

#[derive(Parser, Debug)]
struct Cli {
    /// The address to listen on
    #[arg(short, long, default_value_t = SocketAddr::new(LOCALHOST_V6, 3031))]
    addr: SocketAddr,
    /// The WASM component to load
    component: PathBuf,
}

#[tonic::async_trait]
impl Greeter for MyGreeter {
    async fn say_hello(
        &self,
        _request: Request<HelloRequest>,
    ) -> Result<Response<HelloReply>, Status> {

        let reply = async move {
            let mut mycomp = self.c.lock().await;
            mycomp.call().await.map_err(|x| Status::unknown(x.to_string()))
        }
        .await?;

        let ret = hello_world::HelloReply {
            message: reply.message,
        };

        Ok(Response::new(ret))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Cli::parse();

    let c = MyComponent::new(args.component).await?;

    let greeter = MyGreeter {
        c: Arc::new(Mutex::new(c)),
    };

    Server::builder()
        .add_service(GreeterServer::new(greeter))
        .serve(args.addr)
        .await?;

    Ok(())
}
