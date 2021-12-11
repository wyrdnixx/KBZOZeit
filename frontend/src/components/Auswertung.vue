<template>
  <div >    
    <p>
      Auswertung
    </p>   
        <button class="btn btn-info" v-on:click="TestBannerMessage()">Test-BannerMessage</button>
  </div>
</template>

<script>

export default {
  name: 'Auswertung',  
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
          this.$parent.showAlert("Auswertung")
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
