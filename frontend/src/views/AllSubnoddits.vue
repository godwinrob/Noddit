<template>
  <div class="all-subnoddits">
    <div class="shadow-lg p-3 mb-5 bg-gray rounded">
    <div class="header-div"><h1>All Subnoddits</h1></div>
    <b-list-group>
      
        <b-list-group-item v-for="subnoddit in subnodditArray" :key="subnoddit.subnodditId">
          
            <b-link :to="{ name: 'subnoddit', params: { subnoddit: subnoddit.subnodditName } }">
                <div>/n/{{subnoddit.subnodditName}} - {{subnoddit.subnodditDescription}}</div>
            </b-link>
         
        </b-list-group-item>
        
    </b-list-group> 
  </div>
  </div>
</template>

<script>
import auth from "@/auth.js";
import API_Service from "@/service/API_Service.js";

export default {
  data() {
    return {
        user: auth.getUser(),
        subnodditArray: []
    };
  },
  created() {
      API_Service.getAllSubnoddits().then(
        parsedSubnoddits => (this.subnodditArray = parsedSubnoddits)
      );
  }
};
</script>

<style scoped>
ul {
 list-style-type: none;
}
.header-div {
  margin: 5px;
  color:rgb(112, 112, 112);
}

.li-div {
  max-width: 80%;
  width: 600px;
  padding: 3px;
  border: 3px solid rgb(190, 189, 189);
  margin: 5px;
  border-radius: 4px;
  background-color: rgb(216, 211, 211);
}
</style>