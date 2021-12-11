<template>
  <div id="app">
    <img alt="kbzo logo" src="./assets/logo.png" width="25%" height="25%">
    <Main msg="Stempeln"/>
  </div>
</template>

<script>

import Main from './components/Main.vue'

export default {
  name: 'App',
  components: {
    Main
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #bba3a3;
  margin-top: 60px;
  background-color:rgb(41, 41, 41) ;
}
</style>
