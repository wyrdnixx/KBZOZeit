<template>
  <div class="mainPage">
    <h1>{{ msg }}</h1>        
    <ConnectionStatus :readyState="socketReadyState"/>
    <HelloWorld :socket="this.socket" :message="messageFromServer"  msg="HelloWorldComponent"/>
  </div>
</template>

<script>
import ConnectionStatus from './ConnectionStatus.vue';
import HelloWorld from './HelloWorld.vue'

export default {
  name: 'RootComponent',
  props: {
    msg: String
  },
  components: {
    HelloWorld,
    ConnectionStatus
  },

  data() {
    return {    
      messageFromServer: {}, // To hold the message received from WebSocket
      socket: null, // WebSocket instance
      socketReadyState: 0, // initialer Wert
    };
  },
  mounted() {
    this.setupWebSocket();
  },
  methods: { 
    setupWebSocket() {
      // Create a WebSocket connection
      this.socket = new WebSocket('ws://127.0.0.1:3000/ws');

      this.socket.onopen = () => {
        console.log('WebSocket connection established');
        this.socketReadyState = this.socket.readyState;
            };

      /* this.socket.onmessage = (event) => {
        // Handle the message from the server (echo message)
        console.log('Received from server:', event.data);
        
        this.messageFromServer = event.data; // Update the message
      }; */

      this.socket.onmessage = (event) => {
        console.log('Received from server:', event.data);

        // Parse the string message into a JavaScript object
        try {
          this.messageFromServer = JSON.parse(event.data);
        } catch (error) {
          console.error('Error parsing JSON:', error);
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.socketReadyState = this.socket.readyState;      
        
      };

      this.socket.onclose = () => {
        console.log('WebSocket connection closed');
        this.socketReadyState = this.socket.readyState;    
        this.reconnectSocket();
      };
    },
    async reconnectSocket() {
      this.message = 'Waiting for 5 seconds...';      
      await this.wait(5000); // Wait for 5 seconds
      console.log("trying to reconnect")
      this.setupWebSocket();
    },
    wait(ms) {
      return new Promise(resolve => setTimeout(resolve, ms));
    },
    sendEchoMessage() {
      const echoMessage = {
        type: 'echoTest',
        content: 'Hello from Vue!',
      };

      // Send the JSON message to the server
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify(echoMessage));
        console.log('Sent message:', JSON.stringify(echoMessage));
      } else {
        console.error('WebSocket is not open');
      }
    },
   
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
