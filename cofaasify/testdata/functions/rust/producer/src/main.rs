use std::time::Instant;

use hello_world::{greeter_server::{Greeter, GreeterServer}, HelloReply, HelloRequest};
use prodcon::{producer_consumer_client::ProducerConsumerClient, ConsumeByteRequest};
use tonic::{transport::{Server, Channel, Endpoint}, Request, Response, Status};


pub mod hello_world {
    tonic::include_proto!("helloworld"); // The string specified here must match the proto package name
}

pub mod prodcon {
    tonic::include_proto!("prodcon"); // The string specified here must match the proto package name
}

struct MyGreeter {
    repeats: i32,
    payload: Vec<u8>,
    debug: bool,
    record_inner_lat: bool,
}

#[tonic::async_trait]
impl Greeter for MyGreeter {
    async fn say_hello(
        &self,
        _request: Request<HelloRequest>,
    ) -> Result<Response<HelloReply>, Status> {

        if self.debug {
            println!("Serving request");
        }

        let mut client = ProducerConsumerClient::connect("http://consumer:3030").await
            .map_err(|x| Status::unknown(x.to_string()))?;

        if self.debug {
            println!("Connected to consumer");
        }

        if self.record_inner_lat {
            let start = Instant::now();
            for _ in 1..self.repeats {
                let req = ConsumeByteRequest{
                    value: self.payload.clone()
                };

                let res = client.consume_byte(req).await?.into_inner();
                if self.debug {
                    println!("Performed request that returned with {} and value {}", res.value, res.length)
                }
            }

            let latency = start.elapsed().as_micros() / (self.repeats as u128);

            if self.debug {
                println!("Returing latency {}", latency);
            }

            let ret = HelloReply {
                message: latency.to_string() //reply.message,
            };

            return Ok(Response::new(ret))

        } else {

            for _ in 1..self.repeats {
                let req = ConsumeByteRequest{
                    value: self.payload.clone()
                };

                let res = client.consume_byte(req).await?.into_inner();
                if self.debug {
                    println!("Performed request that returned with {} and value {}", res.value, res.length)
                }
            }

            let ret = HelloReply {
                message: "".into() //reply.message,
            };

            return Ok(Response::new(ret));
        }
    }
}

fn get_env_or_default(env: &str, default: &str) -> String {
    match std::env::var(env) {
        Ok(val) => val,
        Err(_) => default.to_string(),
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let repeats: i32 = get_env_or_default("REPEATS", "1").parse()?;
    let is: isize = get_env_or_default("TRANSFER_SIZE_KB", "1").parse()?;
    let record_inner_lat: bool = get_env_or_default("MEASURE_LAT", "false") == "true";
    let verbose: bool = get_env_or_default("VERBOSE", "false") == "true";

    let us = if is < 0 {
        return Err("Transfer size must be a positive number".into())
    } else {
        is as usize * 1024
    };

    println!("Starting rust server");
    println!("Serving request {} times", repeats);
    println!("with {} KB per request", is);
    if record_inner_lat {
        println!("Running in latency recording mode")
    }

    let greeter = MyGreeter {
        repeats,
        debug: verbose,
        payload: vec![0; us],
        record_inner_lat,
    };

    Server::builder()
        .add_service(GreeterServer::new(greeter))
        //.serve("[::]:3031".parse()?)
        .serve_with_shutdown("[::]:3031".parse()?, async {
            tokio::spawn(tokio::signal::ctrl_c()).await;
        })
        .await?;

    Ok(())
}
