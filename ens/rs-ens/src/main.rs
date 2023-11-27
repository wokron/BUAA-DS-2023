mod args;

#[allow(dead_code)]
mod enslib;

use args::Args;
use clap::Parser;
use enslib::{
    client::{publish, subscribe},
    conn::recv_ens_msg,
};
use log::{error, info};
use std::net::TcpStream;

fn main() {
    env_logger::init();

    let args = Args::parse();

    let addr: String = args.server + ":" + &args.port.to_string();

    if let Ok(mut stream) = TcpStream::connect(addr) {
        match args.command {
            args::Commands::Publish { event } => {
                handle_publish(&mut stream, event);
            }
            args::Commands::Subscribe { topic } => {
                handle_subscribe(&mut stream, topic);
            }
        }
    } else {
        error!("Couldn't connect to server");
    }
}

fn handle_subscribe(stream: &mut TcpStream, topics: Vec<String>) {
    for topic in &topics {
        info!("Subscribe topic \"{}\"", topic);
        if let Err(err) = subscribe(stream, topic) {
            error!("Fail to subscribe topic \"{}\", {}", topic, err);
        }
    }

    loop {
        if let Ok(msg) = recv_ens_msg(stream) {
            info!(
                "Receive event, topic: \"{}\", message: \"{}\".",
                msg.topic, msg.message
            );
        } else {
            error!("error when receiving")
        }
    }
}

fn handle_publish(stream: &mut TcpStream, events: Vec<(String, String)>) {
    for event in &events {
        let topic = &event.0;
        let message = &event.1;

        if let Err(err) = publish(stream, topic, message) {
            error!(
                "Fail to publish event on topic \"{}\", message: \"{}\". {}.",
                topic, message, err
            );
        }
    }
    info!("Publishing events finished.")
}
