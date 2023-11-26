use clap::{Parser, Subcommand};
use std::error::Error;

#[derive(Parser)]
pub struct Args {
    #[arg(short, long, default_value_t = String::from("localhost"))]
    pub server: String,

    #[arg(short, long, default_value_t = 4567)]
    pub port: i32,

    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand)]
pub enum Commands {
    Subscribe {
        #[arg(short, long)]
        topic: Vec<String>,
    },
    Publish {
        #[arg(short, long, value_parser = parse_key_val::<String,String>)]
        event: Vec<(String, String)>,
    },
}

fn parse_key_val<T, U>(s: &str) -> Result<(T, U), Box<dyn Error + Send + Sync + 'static>>
where
    T: std::str::FromStr,
    T::Err: Error + Send + Sync + 'static,
    U: std::str::FromStr,
    U::Err: Error + Send + Sync + 'static,
{
    let pos = s
        .find(':')
        .ok_or_else(|| format!("invalid KEY=value: no `:` found in `{s}`"))?;
    Ok((s[..pos].parse()?, s[pos + 1..].parse()?))
}
