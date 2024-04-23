mod host;
use host::MyComponent;
use tokio::sync::oneshot::{self, Sender};
use tonic::{transport::Server, Request, Response, Status};

use hello_world::greeter_server::{Greeter, GreeterServer};
use hello_world::{HelloReply, HelloRequest};

use clap::Parser;
use futures::lock::Mutex;
use std::net::IpAddr;
use std::net::Ipv6Addr;
use std::net::SocketAddr;
use std::path::PathBuf;
use std::sync::Arc;
use futures::FutureExt;

use std::sync::atomic::AtomicI32;

pub mod hello_world {
    tonic::include_proto!("helloworld"); // The string specified here must match the proto package name
}

pub struct MyGreeter {
    c: Arc<Mutex<MyComponent>>,
    tx: Arc<Mutex<Option<Sender<()>>>>,
    req_count: AtomicI32,
    stop_after: i32,
}

const LOCALHOST_V6: IpAddr = IpAddr::V6(Ipv6Addr::new(0, 0, 0, 0, 0, 0, 0, 1));

#[derive(Parser, Debug)]
struct Cli {
    /// The address to listen on
    #[arg(short, long, default_value_t = SocketAddr::new(LOCALHOST_V6, 3031))]
    addr: SocketAddr,
    /// Enable wasm jitdump profiling
    #[arg(short, long, default_value_t = false)]
    profile: bool,
    /// Stop server after a number of requests (only active when
    /// profiling is enabled)
    #[arg(short, long, default_value_t = 1000)]
    stop_after: i32,
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
            mycomp
                .call()
                .await
                .map_err(|x| Status::unknown(x.to_string()))
        }
        .await?;

        let ret = hello_world::HelloReply {
            message: reply.message,
        };

        let cur_req_count = self.req_count.fetch_add(1, std::sync::atomic::Ordering::Relaxed);
        if cur_req_count >= self.stop_after {
            let mut ch = self.tx.lock().await;
            if ch.is_some() {
                let _ = ch.take().unwrap().send(());
            }
        }

        Ok(Response::new(ret))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Cli::parse();

    let c = MyComponent::new(args.component, args.profile).await?;

    let (tx, rx) = oneshot::channel::<()>();

    let tx_wrap = if args.profile {
        Arc::new(Mutex::new(Some(tx)))
    } else {
        Arc::new(Mutex::new(None))
    };

    let greeter = MyGreeter {
        c: Arc::new(Mutex::new(c)),
        tx: tx_wrap,
        req_count: AtomicI32::new(0),
        stop_after: args.stop_after,
    };

    if args.profile {
        Server::builder()
            .add_service(GreeterServer::new(greeter))
            .serve_with_shutdown(args.addr, rx.map(drop))
            .await?;
    } else {
        Server::builder()
            .add_service(GreeterServer::new(greeter))
            .serve(args.addr)
            .await?;
    }
    Ok(())
}
