# TCP vs UDP
- Both TCP (Transmission Control Protocol) and UDP (User Datagram Protocol) are both Transport Layer Protocols. 
- TCP - connnection-oriented, handshake, reliable, order, congestion/flow controler, slower, heavier.
- UDP - connectionless, no handshake, unreliable, unordered, no congestion\flow controler, faster/lighter.

## Byte stream vs datagrams
- TCP: byte stream (no message boundaries). App must frame messages (length-prefix or delimiter).
- UDP: message/datagram oriented (boundaries preserved). One recv = one datagram (or it's dropped).

## When to use 
- TCP: web (HTTP/HTTPS), file transfer, email, anything requiring correctness/order.
- UDP: real-time voice/video, gaming DNS - Latency-sensitive, can tolerate/handle loss. 


## Framing strategies
- TCP
	- Length-prefix: fixed-size header (e.g, 4-8 bytes) giving payload length; then read n bytes. 
	- Delimiter: e,g, newline; read until delimiter.
- UDP
	- Prefer one logical message per datagram user MTU (~1200 bytes safe). 
	- If splitting is necessary, add app-level header (msg ID, seq, total,), reassamble with timeouts. 

## Connections
- Connection-oritented (TCP)
	- Requires a setup pharse before data: 3-way handshake.
	- Handshake
		1. Client -> Server SYN (A TCP Flag bit) with client ISN( Initial Sequence Number).
		2. Server -> Client: SYM + ACK (ACK = client ISN + 1) with server ISN. 
		3. Client -> Server: ACK (ACK = server ISN + 1).
	- After this, both sides have synced sequence number and can exchange data reliably. 
	- Server must be listening/accepting; connection holds state until closed. 
- Connectionless (UDP)
	- No handshake, no connection state.
	- Any host can send datagram to IP:port at anytime. 
	- Receiver just needs a sock bounded to the port; otherwise packets are dropped. 
	- Each datagram is indepentdent, no build-in reliability, ordering, or congestion control.	  
	
## Key Takeaways
- TCP may split/combines writes, we must parse to reconstruct messages. 
- UDP preserves, each message as sent, but we must handle loss, duplication, and reordering. 