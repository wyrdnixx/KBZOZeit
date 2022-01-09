<template>
  <div >    
    
      <h2> Stempeln </h2>
    <!--  Debug: {{this.sendMessage}}
    
        <br>
        <button class="btn btn-secondary" v-on:click="GetOpenTimer()">Test-offene Timer abfragen</button>
        <br>
        -->
        <div>
          <div v-if="this.openTimerStartTime == '' ">
            <h3>keine offene Z&auml;hlung vorhanden</h3>
            <button class="btn btn-success" v-on:click="Einstempeln()">Einstempeln</button>
          </div>
          <div v-else>
            <h3> Eingestempelt seit:  {{ this.openTimerStartTime | moment("DD.MM.YYYY hh:mm") }}</h3>            
            <!-- <span>Now : {{ (new Date()) | moment("DD.MM.YYYY hh:mm:ss") }}</span>                     -->
            
            <h3>Zeit vergangen:  {{this.TimeDiff}} Stunden</h3>
            <br>
            <button class="btn btn-danger" v-on:click="AusstempelnButton()">Ausstempeln</button>
            <br>
          </div>
        </div>
        <!-- Test <button class="btn btn-info" v-on:click="getTimeDiff()">getTimeDiff</button> -->

  </div>
</template>

<script>
import axios from 'axios';

//const apiURL = window.location.protocol + "//"+ window.location.href +"/api"
//const apiURL = window.location.href +"api"

export default {
  name: 'Stempeln',  
  props: {
    
  },
  data() {
      return {
          sendMessage: {
            MsgType: "TimeAccounting",            
            Name: "",
            Typ:"",
            FromDate:""            
            },
          openTimerStartTime:"",
          TimeDiff:"", 
          
      }
  },
  created() {
    this.sendMessage.Name = this.$parent.username  
    this.GetOpenTimer()

      //refresh - evtl noch verschieben - nur wenn offener timer gefunden wurde
      window.setInterval(() => {
        this.GetOpenTimer()
        this.getTimeDiff()
      }, 10000) // every 10 Seconds
  },
  computed: {    

  },
  methods: {
    getTimeDiff(){
      //var diff =( (new Date(this.openTimerStartTime)) -(new Date())) / 1000;
      var now = new Date()
      var then = new Date(this.openTimerStartTime)
      //console.log("now: " + now)
      //console.log("then: " + then)
      var diff =( now - then)  / 1000;
      //console.log("diff " + diff)
      diff /= (60 * 60);  // dif in stunden umgerechnet
      //console.log("diff " + diff)
      //console.log("diff " + diff.toFixed(2) )     

      this.TimeDiff =diff.toFixed(2) 
    },     
      async GetOpenTimer() {
        console.log("GetOpenTimer: " + this.$parent.APIURL + '/TestApi')
        this.sendMessage.Typ = "getOpenTimer"

        axios.post (this.$parent.APIURL + '/TestApi', this.sendMessage)
        .then ((res) => {
          console.log("Result: " + res.data.FromDate.String)
          this.openTimerStartTime = res.data.FromDate.String
        })
        .catch((error) => {                     
          console.log("Error:"+ error.response.data.Result)
          this.$parent.showAlert(error.response.data)
        })
        .finally(() => {
          //Perform action in always
          this.getTimeDiff()
          
        });
      },

      async Einstempeln() {
          console.log("Einstempeln")            
          
          this.sendMessage.Typ = "Einstempeln"
          this.sendMessage.FromDate = this.openTimerStartTime
             axios.post(this.$parent.APIURL+ '/TestApi', this.sendMessage)
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
                     this.GetOpenTimer()
                 });
    },
    async AusstempelnButton() {
      console.log("Ausstempeln")
            
          // start or stop 
          this.sendMessage.Typ = "Ausstempeln"
          this.sendMessage.FromDate = this.openTimerStartTime          
             axios.post(this.$parent.APIURL+ '/TestApi', this.sendMessage)
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
                     this.GetOpenTimer()
                 });
      
    },   
   
  
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
