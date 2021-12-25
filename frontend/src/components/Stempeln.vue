<template>
  <div >    
    <p>
      Stempeln
      Debu: {{this.sendMessage}}
    </p>   
        <button class="btn btn-info" v-on:click="TestButton()">TestButton</button>
  </div>
</template>

<script>
import axios from 'axios';
const apiURL = window.location.protocol + "//"+ window.location.hostname +":8081/api"

export default {
  name: 'Stempeln',  
  props: {
    
  },
  data() {
      return {
          sendMessage: {
            MsgType: "TimeAccounting",            
            Name: "",
            Typ:""
            }
      }
  },
  created() {
    this.sendMessage.Name = this.$parent.username          
  },
  methods: {
      TestButton() {
          //this.$parent.showAlert("Stempeln")
          //this.TimeAccounting()
          this.GetOpenTimeaccounting()
      },
      async GetOpenTimeaccounting(){
          console.log("GetOpenTimeaccounting")
          this.sendMessage.MsgType= "GetOpenTimeaccounting"
          this.sendMessage.Name= "testuser"
          axios.post(apiURL+'/TestApi',this.sendMessage)
          .then((res) => {
            console.log("Result: "+ res.data.Result)
          })
          .catch((error) => {                     
               console.log("Error:"+ error.response.data.Result)
               this.$parent.showAlert(error.response.data)
          })
           .finally(() => {
               //Perform action in always
           });
      },
      async TimeAccounting() {
          console.log("TimeAccounting")
            
          // start or stop 
          this.sendMessage.Typ = "startAccounting"

             axios.post(apiURL+ '/TestApi', this.sendMessage)
                 .then((res) => {
                     //Perform Success Action
                     console.log("Resut: "+ res.data.Result)
                     if (res.data.Result == "login successfully") {
                       this.$cookies.set("username",this.sendMessage.Name);
                       this.$parent.checkCookie()
                     }
                 })
                .catch((error) => {                     
                     console.log("Error:"+ error.response.data.Result)
                     this.$parent.showAlert(error.response.data)
                })
                 .finally(() => {
                     //Perform action in always
                 });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
