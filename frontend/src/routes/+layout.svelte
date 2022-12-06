<script>
	let socket = new WebSocket("ws://localhost:6969/ws/");
	let messages = [];

	socket.onopen = function(e) {
		alert("[open] Connection established");
	};

	socket.onmessage = function(event) {
		messages.push(event.data)
	};

	socket.onclose = function(event) {
		if (event.wasClean) {
			alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
		} else {
			// e.g. server process killed or network down
			// event.code is usually 1006 in this case
			alert(event.code)
			alert('[close] Connection died');
		}
	};

	socket.onerror = function(error) {
		alert(`[error]`);
	};
</script>

<h1>Socket: {socket.readyState}</h1>
<h1>Messages</h1>
{#each messages as message}
	<p>{message}</p>
{/each}
