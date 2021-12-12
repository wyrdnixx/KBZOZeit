<template>
  <div >    
    <p>
      Nacherfassen
    </p>
        <div>
          <table class="table table-dark">
  <thead>
    <tr>      
      <th scope="col">Datum</th>
      <th scope="col">Von</th>
      <th scope="col">Bis</th>
    </tr>
  </thead>
  <tbody>
    <tr>   
      <td><input type="date" id="entryDate" /></td>
      <td><input type="time" id="entryFrom" /></td>
      <td><input type="time" id="entrytimeTo" /></td>               
    </tr>
  </tbody>
</table>
<input type="text" id="entryText" placeholder="was wurde gemacht?"/>
        </div>
      <hr>
    <button class="btn btn-info" v-on:click="TestBannerMessage()">Nacherfassen</button>
  </div>
</template>

<script>

export default {
  name: 'Nacherfassen',  
  props: {
    
  },
  data() {
      return {

      }
  },
  created() {

  },
  methods: {
      TestBannerMessage() {
          this.$parent.showAlert("Nacherfassen")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
