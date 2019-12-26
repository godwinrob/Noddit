<template>
  <div class="profile-edit">
      <div class="shadow p-3 mb-5 bg-gray rounded">
    <b-card
      no-body
      style="max-width: 30rem;"
      :img-src="require('@/assets/img/Noddit-Logo.png')"
      img-alt="Image"
      img-top
    >
      <template v-slot:header>
        <h4 class="mb-0">{{userName.sub}}/{{userName.rol}}</h4>
      </template>


      <b-card-body class="card-body">
        <b-card-title>Account Settings</b-card-title>
        <b-card-sub-title class="mb-2"></b-card-sub-title>
        <b-card-text></b-card-text>
      </b-card-body>

      <div role="tablist">
        <b-card no-body class="mb-1">
         <!-- <gb-card-header header-tag="header" class="p-1" role="tab">-->
            <gb-button  block href="#" v-b-toggle.accordion-1 color="black" >Change Username</gb-button>
          <!--</gb-card-header>-->
          <b-collapse id="accordion-1" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <b-card-text></b-card-text>
              <b-card-text>
                <div>
                  <b-form-input v-model="user.newUsername" placeholder="Update Username"></b-form-input>
                </div>
                <br>
                <div>
                    <button v-b-modal.modal-1 block class="btn-block btn-dark" @click.prevent="changeUsername" >Submit</button>

                    <b-modal id="modal-1" title="Noddit">
                         <p class="my-4">Your Username has been updated!</p>
                    </b-modal>
                </div>
              </b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card>

        <b-card no-body class="mb-1">
         
            <!--<gb-button block href="#" v-b-toggle.accordion-2 color="black" >Update Name</gb-button>
         
          <b-collapse id="accordion-2" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <b-card-text>
                <div>
                  <b-form-input v-model="updateFirstName" placeholder="Update User First name"></b-form-input>
                  <br>
                  <b-form-input v-model="updateLastName" placeholder="Update User Last name"></b-form-input>
                </div>
                <br>
                 <div>
                    <gb-button v-b-modal.modal-2 block color="black">Submit</gb-button>

                    <b-modal id="modal-2" title="Noddit">
                         <p class="my-4">Your first and last name has been updated!</p>
                    </b-modal>
                </div>
              </b-card-text>
            </b-card-body>
          </b-collapse>-->
        </b-card>

        <b-card no-body class="mb-1">
          <!--<b-card-header header-tag="header" class="p-1" role="tab">-->
            <gb-button block href="#" v-b-toggle.accordion-3 color="black">Change Password</gb-button>
          <!--</b-card-header>-->
          <b-collapse id="accordion-3" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <b-card-text>
                <div>
                  <b-form-input placeholder="Old Password"></b-form-input>
                  <br>
                  <b-form-input placeholder="New Password"></b-form-input>
                  <br>
                  <b-form-input placeholder="Re-enter Password"></b-form-input>
                </div>
                <br>
                 <div>
                    <button v-b-modal.modal-3 block class="btn-block btn-dark">Submit</button>

                    <b-modal id="modal-3" title="Noddit">
                         <p class="my-4">Your Password has been updated!</p>
                    </b-modal>
                </div>
              </b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card>

          <b-card no-body class="mb-1">
          <!--<b-card-header header-tag="header" class="p-1" role="tab">-->
            <gb-button block href="#" v-b-toggle.accordion-4 color="black">Change E-mail</gb-button>
         <!-- </b-card-header>-->
          <b-collapse id="accordion-4" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <b-card-text>
                <div>
                  
                  <b-form-input v-model="user.emailAddress" placeholder="Update E-mail"></b-form-input>
                </div>
                <br>
                <div>
                    <button v-b-modal.modal-4 block class="btn-block btn-dark">Submit</button>

                    <b-modal id="modal-4" title="Noddit">
                         <p class="my-4">Your E-mail has been updated!</p>
                    </b-modal>
                </div>
              </b-card-text>
            </b-card-body>
          </b-collapse>
        </b-card>

          <b-card no-body class="mb-1">
     
           <!-- <gb-button block href="#" v-b-toggle.accordion-5 color="black">Change Avatar</gb-button>
         
          <b-collapse id="accordion-5" accordion="my-accordion" role="tabpanel">
            <b-card-body>
              <b-card-text>
                <div>
                  <b-form-input v-model="updateAvatar" placeholder=" Avatar URL"></b-form-input>
                </div>
                <br>
                <div>
                    <gb-button v-b-modal.modal-5 block color="black">Submit</gb-button>

                    <b-modal id="modal-5" title="BootstrapVue">
                         <p class="my-4">Your Avatar has been updated!</p>
                    </b-modal>
                </div>
              </b-card-text>
            </b-card-body>
          </b-collapse>-->
        </b-card>
      </div>

      <b-card-body>
        <a href="#" class="card-link">Contact Us</a>
        <a href="#" class="card-link">Saved Noddits</a>
      </b-card-body>
    </b-card>
  </div>
  </div>
</template>

<script>
import auth from "../auth";
import API_Service from "@/service/API_Service.js";

export default {
  data() {
    return {
      userName: auth.getUser(),
      updateFirstName:'',
      updateLastName:'',
      updateEmail:'',
      updateAvatar:'',
      user: {
        newUsername: '',
        emailAddress: '',
        username: auth.getUser().sub
      }
    };
  },
  methods: {
    changeUsername() {
      API_Service.updateUsername(this.user);
      auth.logout();
      this.$router.push("/login");
    }
  }
};
</script>

<style>
.card-body{
  background-color: rgb(204, 202, 202);
}


.profile-edit {
  
  text-align: center;
}
</style>