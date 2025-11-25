# hanashi

Peer-to-peer chat written in Go

## Instructions

### Clone this repository

```bash
https://github.com/FhmiSddq/ProyekJarkom.git

```

#### (Optional) Configure the app

Modify the `.env` file

### Running the app

Run the app first with `go run app/main.go`

#### Connecting to another device

```bash
/connect 10.0.0.3:6660
```

#### Sending messages

>> message will be sent to all connected clients

##### On the sending side

```bash
testing
```

##### On the receiving side(s)

```bash
2025/11/19 13:18:09 [10.0.0.2:35532]: testing
```

#### Exit or disconnect

To exit, run

```bash
/exit
```

To disconnect, run

```bash
/disconnect 10.0.0.2:6660
2025/11/19 06:38:58 stopped handling connection: 10.0.0.2:6660
2025/11/19 06:38:58 read tcp 10.0.0.2:53916->100.84.161.27:6660: use of closed network connection
2025/11/19 06:38:58 stopped handling connection: 10.0.0.2:6660
```

