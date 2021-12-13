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
        Wer sind Sie? <input type="text" placeholder="Name" v-model="Ident.Name"/>
        <hr>
        Debug -> Name: {{this.Ident.Name}} UUID: {{this.Ident.UUID}} Status: {{this.Ident.Status}}
        <button class="btn btn-secondary" v-on:click="RegisterIdent()">Register</button>
    </div>
  </div>
</template>

<script>
import { uuid } from 'vue-uuid'; 
import axios from 'axios';
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Register', 
  components:{
    
  },
  props: {
    
  },
  data() {
      return {
          Ident: {
            Name: "",
            UUID: "",
            Status:""
          } 
      }
  },
  created() {

  },
  methods: {

    testMethod() {
      this.Ident.UUID = uuid.v4()
      this.Ident.Status=""       
    },
    async RegisterIdent() {
          console.log("RegisterIdent")
          if (this.Ident.Name === "") {
            this.$parent.showAlert("Kein Username angegeben!")
          } else {
                  
            await   axios.post(apiURL + "/RegisterIdent", {           
              Ident: this.Ident
            }, {
              headers: {
                'Content-Type': 'application/json',
              }
            }).catch((error)=> this.$parent.showAlert("Server returned an Error:\n" + error));  
    
          }
    },        
  }       
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
