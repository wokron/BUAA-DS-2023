use super::msg::{ENSMsg, ENS_MGS_SIZE};
use anyhow::Result;
use std::{
    io::{Read, Write},
    net::TcpStream,
};

pub fn send_ens_msg(stream: &mut TcpStream, msg: &ENSMsg) -> Result<()> {
    let data = msg.encode();

    stream.write_all(data.as_slice())?;

    return Ok(());
}

pub fn recv_ens_msg(stream: &mut TcpStream) -> Result<ENSMsg> {
    let mut recv_data = [0u8; ENS_MGS_SIZE];

    stream.read_exact(&mut recv_data)?;

    return ENSMsg::decode(&recv_data.to_vec());
}
