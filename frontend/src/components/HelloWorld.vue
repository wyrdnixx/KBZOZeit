<template>
  <div>
    <h1>This is: {{ msg }}</h1>
    <button @click="this.$parent.sendEchoMessage">Send Echo Message</button>
    <div>    
      <h1 v-if="message && message.content">Message from WebSocket: {{ message.content }}</h1>
    </div>
  </div>
  
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String,
    socket: {
      type: Object,
      required: true, // Expecting a WebSocket instance from App.vue
    },
    message: {
      type: Object,
      required: true, // Message received from the WebSocket server (now an object)
    },
  },
  data() {
    return {    
      //messageFromServer: '', // To hold the message received from WebSocket      
    };
  },
  methods: {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
