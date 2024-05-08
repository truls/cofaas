mod description;
mod executor;

use anyhow::Result;
use executor::Executor;
use std::path::PathBuf;

use clap::Parser;

#[derive(Parser)]
///Transforms a FaaS application function graph into a single
/// monolithic application using the CoFaaS methodoloty
struct Cli {
    /// The yaml description of the FaaS application
    description_file: PathBuf,
}

fn main() -> Result<()> {
    let args = Cli::parse();

    let ex = Executor::new(args.description_file)?;
    ex.execute()
}
