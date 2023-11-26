use std::net::TcpStream;

use super::{
    conn::send_ens_msg,
    msg::{ENSMsg, ENSMsgType},
};
use anyhow::{Ok, Result};

pub fn publish(stream: &mut TcpStream, topic: &String, message: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::PUBLISH as u8,
        topic: topic.to_string(),
        message: message.to_string(),
    };

    send_ens_msg(stream, &msg)?;
    Ok(())
}

pub fn subscribe(stream: &mut TcpStream, topic: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::SUBSCRIBE as u8,
        topic: topic.to_string(),
        message: String::new(),
    };

    send_ens_msg(stream, &msg)?;
    Ok(())
}

pub fn unsubscribe(stream: &mut TcpStream, topic: &String) -> Result<()> {
    let msg = ENSMsg {
        msg_type: ENSMsgType::UNSUBSCRIBE as u8,
        topic: topic.to_string(),
        message: String::new(),
    };

    send_ens_msg(stream, &msg)?;
    Ok(())
}
