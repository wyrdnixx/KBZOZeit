<template>
  <div>

    <h1 v-if="this.UserAuthenticated">{{ msg }} Wilkommen {{this.username}}</h1>
    <h1 v-else> {{ msg }} - Bitte anmelden </h1>
    <p>
      Main App
    </p>

    <!-- Toast Message banner  -->

    <b-alert
      :show="AlertDismissCountDown"
      dismissible
      variant="warning"
      @dismissed="AlertDismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      <p> {{this.AlertWebsiteAlertMessage}} </p>
        <!-- {{ AlertDismissCountDown }} </p> -->
      <b-progress
        variant="warning"
        :max="AlertDismissSecs"
        :value="AlertDismissCountDown"
        height="4px"
      ></b-progress>
    </b-alert>

<!-- Toast Message banner  -->

  
    <button class="btn btn-secondary" v-on:click="TestChangeAuth()">Clear Auth</button>
    <button class="btn btn-warning" v-on:click="NavChange('Admin')">Admin</button>
    <div v-if="!this.UserAuthenticated">
      <Login /> 

    </div>
    <div v-else>
          <button class="btn btn-info" v-on:click="NavChange('s')">Stempeln</button>
          <button class="btn btn-info" v-on:click="NavChange('n')">Nacherfassen</button>
          <button class="btn btn-info" v-on:click="NavChange('a')">Auswertung</button>
    <div v-if="this.NavCurrentSelected == 's'" >
      <Stempeln />
    </div>
    <div v-if="this.NavCurrentSelected == 'n'" >
      <Nacherfassen />
    </div>
    <div v-if="this.NavCurrentSelected == 'a'" >
      <Auswertung />
    </div>
      
    </div>
    <div v-if="this.NavCurrentSelected == 'Admin'">
        <Admin />
    </div>
  </div>
</template>

<script>
import Login from './Login.vue'
import Stempeln from './Stempeln.vue'
import Nacherfassen from './Nacherfassen.vue'
import Auswertung from './Auswertung.vue'
import Admin from './Admin.vue'


export default {
  name: 'Main',
  components: {
    Login,
    Stempeln,
    Auswertung,
    Nacherfassen,
    Admin

  },
  props: {
    msg: String
  },
  data() {
    return {
      Username:"",
      UserAuthenticated: false,
      AlertDismissSecs: 4,
      AlertDismissCountDown: 0,
      AlertShowDismissibleAlert: false,
      AlertWebsiteAlertMessage: "",
      NavCurrentSelected: "s"   
    }
    
  },
  created() {
    
    //this.UserAuthenticated = false
    this.checkCookie()
  },
  methods: {

    checkCookie() {
            this.username = this.$cookies.get("username");
            console.log("Cookie got: " + this.username)
            if (this.username != null) {
              this.UserAuthenticated = true
            } else {
              this.UserAuthenticated = false
            }
    },
    TestChangeAuth() {
      console.log("TestChangeAuth")
      this.$cookies.remove("username");

      this.checkCookie()
            
      
    },
    countDownChanged(AlertDismissCountDown) {
      this.AlertDismissCountDown = AlertDismissCountDown
    },
    showAlert(msg) {

      this.AlertWebsiteAlertMessage = msg
      this.AlertDismissCountDown = this.AlertDismissSecs
    },

    NavChange(option){
      this.NavCurrentSelected = option
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
