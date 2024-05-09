<template>
  <div class="container-fluid mt-4" style="padding-left: 300px;">
    <div class="row align-items-center">
      <div class="col-lg-4 col-sm-8">
        <!-- <nav-pill /> -->
      </div>
    </div>
    <!-- <div class="row">
      <div class="col-lg-8">
        <div class="row mt-4">
          <div class="col-xl-6 mb-xl-0 mb-4">
            <master-card />
          </div>
          <div class="col-xl-6">
            <div class="row">
              <div class="col-md-6">
                <default-info-card
                  icon="text-white fas fa-landmark"
                  title="Salary"
                  description="Belong Interactive"
                  value="+$2000"
                />
              </div>
              <div class="col-md-6">
                <default-info-card
                  icon="text-white fab fa-paypal"
                  title="Paypal"
                  description="Freelance Payment"
                  value="$455.00"
                />
              </div>
            </div>
          </div>
          <div class="col-md-12 mb-4">
            <payment-methods-card />
          </div>
        </div>
      </div>
      <div class="col-lg-4">
        <invoice-card />
      </div>
    </div> -->
    <div class="row mt-4">
      <div class="col-md-9">
        <billing-card
          title="区块链信息"
          :bills="BlockChainInfo"
        />
      </div>
      <!-- <div class="col-md-5 mt-4">
        <ranking-list-card
          :horizontal-break="false"
          :card="{
            title: 'Transactions',
            date: '23 - 30 March 2021',
            subtitle: 'Newest',
          }"
          :item="[
            {
              title: 'Netflix',
              date: '27 March 2020, at 12:30 PM',
              amount: '- $ 2,500',
              icon: 'fa-arrow-down',
              color: 'danger',
            },
            {
              title: 'Apple',
              date: '23 March 2020, at 04:30 AM',
              amount: '+ $ 2,000',
              icon: 'fa-arrow-up',
              color: 'success',
            },
          ]"
        >
          <ranking-list
            title="Yesterday"
            :item="[
              {
                title: 'Stripe',
                date: '26 March 2020, at 13:45 PM',
                amount: '+ $ 750',
                icon: 'fa-arrow-up',
                color: 'success',
              },
              {
                title: 'HubSpot',
                date: '26 March 2020, at 12:30 PM',
                amount: '+ $ 1,000',
                icon: 'fa-arrow-up',
                color: 'success',
              },
              {
                title: 'Creative Tim',
                date: '26 March 2020, at 08:30 AM',
                amount: '+ $ 2,500',
                icon: 'fa-arrow-up',
                color: 'success',
              },
              {
                title: 'Webflow',
                date: '26 March 2020, at 05:00 AM',
                amount: 'Pending',
                icon: 'fa-info',
                color: 'dark',
              },
            ]"
            :horizontal-break="false"
          />
        </ranking-list-card>
      </div> -->
    </div>
  </div>
</template>

<script>
import setTooltip from "@/assets/js/tooltip.js";
// import NavPill from "./components/NavPill.vue";
// import MasterCard from "@/examples/Cards/MasterCard.vue";
// import DefaultInfoCard from "@/examples/Cards/DefaultInfoCard.vue";
// import PaymentMethodsCard from "./components/PaymentMethodsCard.vue";
// import InvoiceCard from "./components/InvoiceCard.vue";
import BillingCard from "./components/BillingCard.vue";
// import RankingList from "@/examples/Cards/RankingList.vue";
// import RankingListCard from "@/examples/Cards/RankingListCard.vue";




import { onMounted, ref } from 'vue'
import axios from 'axios'

const BlockChainInfo=ref()

const getBlockChainInfo=()=>{
  axios.get('http://localhost:8081/blockchain_info')
  .then(response=>{
    BlockChainInfo.value=response.data

  })
  .catch(err=>{
    console.log(err)
  })

}
// const getBlockChainInfo =() =>{

// }

export default {
  name: "Billing",
  setup(){
    onMounted(()=>{
      getBlockChainInfo();
    }
    );

    return{
      BlockChainInfo
    }

  },
  components: {
    // NavPill,
    // MasterCard,
    // DefaultInfoCard,
    // PaymentMethodsCard,
    // InvoiceCard,
    BillingCard,
    // RankingList,
    // RankingListCard,
  },
  mounted() {
    setTooltip();
    // getBlockChainInfo();
  },
};
</script>
