<template>
  <div class="card h-100">
    <div class="p-3 pb-0 card-header">
      <div class="row">
        <div class="col-md-8 d-flex align-items-center">
          <h6 class="mb-0">{{ title }}</h6>
        <div style="height: 10px;"></div>
        </div>
        <div class="col-md-4 text-end">
          <a :href="action.route">
            <!-- <i
              class="text-sm fas fa-user-edit text-secondary"
              data-bs-toggle="tooltip"
              data-bs-placement="top"
              :title="action.tooltip"
            ></i> -->
          </a>
        </div>
      </div>
    </div>
    <div class="p-3 card-body">
      <!-- <p class="text-sm">
        {{ description }}
      </p> -->
      <!-- <hr class="my-4 horizontal gray-light" /> -->
      <ul class="list-group">
        <li class="pt-0 text-sm border-0 list-group-item ps-0">
          <strong class="text-dark">用户名:</strong> &nbsp;         
          <p></p>
          {{ info.user_name }}
        </li>
        <!-- <div style="height: 50px;"></div> -->
        <li class="text-sm border-0 list-group-item ps-0">
          <strong class="text-dark">钱包地址:</strong> &nbsp; 
          <p></p>
          {{ info.wallet_address }}
        </li>
        <!-- <div style="height: 50px;"></div> -->
        <li class="text-sm border-0 list-group-item ps-0">
          <strong class="text-dark">茅台数量:</strong> &nbsp;
          <p></p>
           {{ info.user_balance }}
        </li>
        <!-- <div style="height: 50px;"></div> -->
        <li class="text-sm border-0 list-group-item ps-0">
          <strong class="text-dark">个人公钥:</strong> &nbsp; 
          <p></p>
          {{ info.user_pub_key }}
        </li>
        <li class="text-sm border-0 list-group-item ps-0">
          <strong class="text-dark">身份组:</strong> &nbsp;
          <p></p>
          {{ info.user_identity }}
        </li>
        <!-- <div style="height: 50px;"></div> -->
        <li class="pb-0 border-0 list-group-item ps-0">
          <!-- <strong class="text-sm text-dark">Social:</strong> &nbsp; -->
          <a
            v-for="({ icon, link }, index) of social"
            :key="index"
            class="py-0 mb-0 btn-simple ps-1 pe-2"
            :href="link"
          >
            <font-awesome-icon :icon="icon" />
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

export default {
  name: "ProfileInfoCard",
  components: {
    FontAwesomeIcon,
  },
  props: {
    title: {
      type: String,
      default: "",
    },
    description: {
      type: String,
      default: "",
    },
    info: {
      type: Object,
      fullName: String,
      mobile: String,
      email: String,
      location: String,
      default: () => {},
    },
    social: {
      type: Array,
      link: String,
      icon: String,
      default: () => [],
    },
    action: {
      type: Object,
      route: String,
      tooltip: String,
      default: () => ({
        route: "javascript:;",
      }),
    },
  },
};
</script>
