package mtproto

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type cmdInfo struct {
	cmd  uint32
	name string
}

var cmds []*cmdInfo
var cmdToInfo = make(map[uint32]*cmdInfo)

func RegisterCmd(cmd uint32, name, def string) {
	cinfo := &cmdInfo{cmd, name}
	cmds = append(cmds, cinfo)
	cmdToInfo[cmd] = cinfo
}

func AddSchema(schema string) {
	for _, line := range strings.Split(schema, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0:2] == "//" {
			continue
		}
		if line[0:3] == "---" {
			continue
		}

		AddSchemaLine(line)
	}
}

func AddSchemaLine(line string) {
	line = strings.TrimSpace(line)

	if strings.HasSuffix(line, ";") {
		line = line[:len(line)-1]
	}

	fields := strings.Fields(line)
	name, cmd := parseCombinatorName(fields[0])
	if cmd != 0 {
		RegisterCmd(cmd, name, line)
	}
}

func DescribeCmd(cmd uint32) string {
	if cmd == 0 {
		return "none"
	} else if cinfo := cmdToInfo[cmd]; cinfo != nil {
		return fmt.Sprintf("%s#%08x", cinfo.name, cmd)
	} else {
		return fmt.Sprintf("#%08x", cmd)
	}
}

func DescribeCmdOfPayload(b []byte) string {
	return DescribeCmd(CmdOfPayload(b))
}

func parseCombinatorName(s string) (string, uint32) {
	idx := strings.IndexRune(s, '#')
	if idx < 0 {
		return s, 0
	}

	name := s[:idx]
	cmdstr := s[idx+1:]
	if len(cmdstr) != 8 {
		log.Panicf("invalid schema, cmd hex code not 8 chars in %#v", s)
	}
	cmd, err := strconv.ParseUint(cmdstr, 16, 32)
	if err != nil {
		log.Panicf("invalid schema, cannot parse hex in %#v: %v", s, err)
	}
	return name, uint32(cmd)
}

func init() {
	AddSchema(mtprotoSchema)
}

const mtprotoSchema = `
05162 ? = Int;
long ? = Long;
double ? = Double;
string ? = String;

vector {t:Type} # [ t ] = Vector t;

int128 4*[ int ] = Int128;
int256 8*[ int ] = Int256;

resPQ#05162463 nonce:int128 server_nonce:int128 pq:bytes server_public_key_fingerprints:Vector<long> = ResPQ;

p_q_inner_data#83c95aec pq:bytes p:bytes q:bytes nonce:int128 server_nonce:int128 new_nonce:int256 = P_Q_inner_data;


server_DH_params_fail#79cb045d nonce:int128 server_nonce:int128 new_nonce_hash:int128 = Server_DH_Params;
server_DH_params_ok#d0e8075c nonce:int128 server_nonce:int128 encrypted_answer:bytes = Server_DH_Params;

server_DH_inner_data#b5890dba nonce:int128 server_nonce:int128 g:int dh_prime:bytes g_a:bytes server_time:int = Server_DH_inner_data;

client_DH_inner_data#6643b654 nonce:int128 server_nonce:int128 retry_id:long g_b:bytes = Client_DH_Inner_Data;

dh_gen_ok#3bcbf734 nonce:int128 server_nonce:int128 new_nonce_hash1:int128 = Set_client_DH_params_answer;
dh_gen_retry#46dc1fb9 nonce:int128 server_nonce:int128 new_nonce_hash2:int128 = Set_client_DH_params_answer;
dh_gen_fail#a69dae02 nonce:int128 server_nonce:int128 new_nonce_hash3:int128 = Set_client_DH_params_answer;

rpc_result#f35c6d01 req_msg_id:long result:Object = RpcResult;
rpc_error#2144ca19 error_code:int error_message:string = RpcError;

rpc_answer_unknown#5e2ad36e = RpcDropAnswer;
rpc_answer_dropped_running#cd78e586 = RpcDropAnswer;
rpc_answer_dropped#a43ad8b7 msg_id:long seq_no:int bytes:int = RpcDropAnswer;

future_salt#0949d9dc valid_since:int valid_until:int salt:long = FutureSalt;
future_salts#ae500895 req_msg_id:long now:int salts:vector<future_salt> = FutureSalts;

pong#347773c5 msg_id:long ping_id:long = Pong;

destroy_session_ok#e22045fc session_id:long = DestroySessionRes;
destroy_session_none#62d350c9 session_id:long = DestroySessionRes;

new_session_created#9ec20908 first_msg_id:long unique_id:long server_salt:long = NewSession;

msg_container#73f1f8dc messages:vector<%Message> = MessageContainer;
message msg_id:long seqno:int bytes:int body:Object = Message;
msg_copy#e06046b2 orig_message:Message = MessageCopy;

gzip_packed#3072cfa1 packed_data:bytes = Object;

msgs_ack#62d6b459 msg_ids:Vector<long> = MsgsAck;

bad_msg_notification#a7eff811 bad_msg_id:long bad_msg_seqno:int error_code:int = BadMsgNotification;
bad_server_salt#edab447b bad_msg_id:long bad_msg_seqno:int error_code:int new_server_salt:long = BadMsgNotification;

msg_resend_req#7d861a08 msg_ids:Vector<long> = MsgResendReq;
msgs_state_req#da69fb52 msg_ids:Vector<long> = MsgsStateReq;
msgs_state_info#04deb57d req_msg_id:long info:bytes = MsgsStateInfo;
msgs_all_info#8cc0d131 msg_ids:Vector<long> info:bytes = MsgsAllInfo;
msg_detailed_info#276d3ec6 msg_id:long answer_msg_id:long bytes:int status:int = MsgDetailedInfo;
msg_new_detailed_info#809db6df answer_msg_id:long bytes:int status:int = MsgDetailedInfo;

---functions---

req_pq#60469778 nonce:int128 = ResPQ;

req_DH_params#d712e4be nonce:int128 server_nonce:int128 p:bytes q:bytes public_key_fingerprint:long encrypted_data:bytes = Server_DH_Params;

set_client_DH_params#f5045f1f nonce:int128 server_nonce:int128 encrypted_data:bytes = Set_client_DH_params_answer;

rpc_drop_answer#58e4a740 req_msg_id:long = RpcDropAnswer;
get_future_salts#b921bd04 num:int = FutureSalts;
ping#7abe77ec ping_id:long = Pong;
ping_delay_disconnect#f3427b8c ping_id:long disconnect_delay:int = Pong;
destroy_session#e7512126 session_id:long = DestroySessionRes;

http_wait#9299359f max_delay:int wait_after:int max_wait:int = HttpWait;
`
