use anyhow::Result;
use serde::Deserialize;
use std::{fs::File, io::Read, path::PathBuf};

#[derive(Deserialize, Debug, PartialEq)]
#[serde(rename_all = "kebab-case")]
struct ProtoFile {
    name: String,
    path: PathBuf,
}

#[derive(Deserialize, Debug, PartialEq)]
#[serde(rename_all = "lowercase")]
enum Language {
    Go,
    JavaScript,
}

#[derive(Deserialize, Debug, PartialEq)]
#[serde(rename_all = "kebab-case")]
struct Function {
    name: String,
    language: Language,
    export: String,
    import: Option<String>,
}

#[derive(Deserialize, Debug, PartialEq)]
#[serde(rename_all = "kebab-case")]
pub struct AppDescription {
    dest_dir: PathBuf,
    proto_files: Vec<ProtoFile>,
    functions: Vec<Function>,
}

impl AppDescription {
    pub fn from_file(file_name: PathBuf) -> Result<Self> {
        let mut f = File::open(file_name)?;
        let mut contents = "".to_string();
        f.read_to_string(&mut contents).unwrap();
        println!("Parsing json: {}", contents);
        let deserialized: AppDescription = serde_yaml::from_str(contents.as_str()).unwrap();
        Ok(deserialized)
    }
}

#[cfg(test)]
mod tests {
    use std::path::PathBuf;

    use crate::description::{AppDescription, Function, Language, ProtoFile};

    #[test]
    fn deserialize() {
        let deserialized =
            AppDescription::from_file(PathBuf::from("test_data/description.yaml")).unwrap();
        let expected = AppDescription {
            dest_dir: PathBuf::from("output"),
            proto_files: vec![
                ProtoFile {
                    name: "helloworld".to_string(),
                    path: PathBuf::from("protos/helloworld.proto"),
                },
                ProtoFile {
                    name: "prodcon".to_string(),
                    path: PathBuf::from("protos/prodcon.proto"),
                },
            ],
            functions: vec![
                Function {
                    name: "producer".to_string(),
                    language: Language::Go,
                    export: "helloworld".to_string(),
                    import: Some("prodcon".to_string()),
                },
                Function {
                    name: "consumer".to_string(),
                    language: Language::Go,
                    export: "prodcon".to_string(),
                    import: None,
                },
            ],
        };

        println!("deserialized = {:?}", deserialized);

        assert_eq!(deserialized, expected);
    }
}
