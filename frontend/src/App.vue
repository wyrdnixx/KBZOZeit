<template>  
  <RootComponent msg="Welcome to RootComponent"/>
  
</template>

<script>
import RootComponent from './components/RootComponent.vue'

export default {
  name: 'App',
  components: {
    RootComponent    
  },
  data() {
    return {
      messageFromServer: '', // Data from WebSocket
    };
  },

};

</script>

<style>
html, body {
  height: 100%; /* Make sure html and body cover the full page */
  margin: 0;
  background-color: rgb(37, 37, 37); /* Set the background of the entire page to black */
}
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;    
    margin-top: 10px;     
  }

</style>
