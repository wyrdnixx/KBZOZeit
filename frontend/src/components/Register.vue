<template>
  <div >    
    <p>
      Register-Tab
    </p>
    <div>  
        Es wurde kein gülltiges UserCookie gefunden.
        <br>
        Bitte registrieren
        <hr>
        Wer sind Sie? <input type="text" placeholder="Name" v-model="UserName"/>
        <hr>
        Debug - Name: {{this.UserName}}
    </div>
  </div>
</template>

<script>

export default {
  name: 'Register', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {
          UserName: ""
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
