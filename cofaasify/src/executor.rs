use std::path::PathBuf;

use anyhow::Result;

use crate::description::AppDescription;

pub struct Executor {
    descr: AppDescription,
}

impl Executor {
    pub fn new(description_file: PathBuf) -> Result<Self> {
        let descr = AppDescription::from_file(description_file)?;
        Ok(Executor { descr })
    }

    pub fn execute(&self) -> Result<()> {
        Ok(())
    }
}
