use prodcon::{
    producer_consumer_server::{ProducerConsumer, ProducerConsumerServer},
    ConsumeByteReply, ConsumeByteRequest,
};
use tonic::{transport::Server, Request, Response, Status};

pub mod prodcon {
    tonic::include_proto!("prodcon"); // The string specified here must match the proto package name
}

struct MyProdCon {
    debug: bool,
}

#[tonic::async_trait]
impl ProducerConsumer for MyProdCon {
    async fn consume_byte(
        &self,
        request: Request<ConsumeByteRequest>,
    ) -> Result<Response<ConsumeByteReply>, Status> {

        if self.debug {
            println!("Handling request");
        }

        let req = request.into_inner();

        let res = ConsumeByteReply {
            value: true,
            length: req.value.len() as i32,
        };

        if self.debug {
            println!("Received a message of {} bytes", res.length);
        }

        Ok(Response::new(res))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Starting rust consumer server");

    let prodcon = MyProdCon { debug: false };

    Server::builder()
        .add_service(ProducerConsumerServer::new(prodcon))
        // .serve("[::]:3030".parse()?)
        .serve_with_shutdown("[::]:3030".parse()?, async {
            tokio::spawn(tokio::signal::ctrl_c()).await;
        })
        .await?;

    Ok(())
}
