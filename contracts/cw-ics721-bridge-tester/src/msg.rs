use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::IbcTimeout;

#[cw_serde]
pub enum AckMode {
    // Messages should respond with an error ACK.
    Error,
    // Messages should respond with a success ACK.
    Success,
}

#[cw_serde]
pub struct InstantiateMsg {
    pub ack_mode: AckMode,
}

#[cw_serde]
pub enum ExecuteMsg {
    CloseChannel {
        channel_id: String,
    },
    SendPacket {
        channel_id: String,
        timeout: IbcTimeout,

        data: cw_ics721_bridge::ibc::NonFungibleTokenPacketData,
    },
    SetAckMode {
        ack_mode: AckMode,
    },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Gets the current ack mode. Returns `AckMode`.
    #[returns(AckMode)]
    AckMode {},
    /// Gets the mode of the last ack this contract received. Errors
    /// if no ACK has ever been received. Returns `AckMode`.
    #[returns(AckMode)]
    LastAck {},
}
