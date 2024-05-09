<template>
  <div class="py-4 container-fluid">
    <div class="row">
      <div style="height: 50px;"></div>
      <div class="col-xl-6 col-sm-6 mb-xl-0 mb-4">
        <mini-statistics-card
          title="原料余量"
          :value="Raw_balance"
          :percentage="{
            value: '+3%',
            color: 'text-success',
          }"
          :icon="{
            component: ' ni ni-world',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div>
      <div class="col-xl-6 col-sm-6 mb-xl-0 mb-4">
        <mini-statistics-card
          title="贵州生产商余量"
          :value="Producer_balance_A"
          :percentage="{
            value: '-2%',
            color: 'text-danger',
          }"
          :icon="{
            component: 'ni ni-paper-diploma',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div>
      <!-- <div class="col-xl-4 col-sm-6 mb-xl-0">
        <mini-statistics-card
          title="重庆生产商余量"
          :value="Producer_balance_B"
          :percentage="{
            value: '+5%',
            color: 'text-success',
          }"
          :icon="{
            component: 'ni ni-cart',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div> -->
    </div>

    <div class="row">
      <div class="col-xl-4 col-sm-6 mb-xl-0 mb-4">
        <div style="height: 30px;"></div>
        <mini-statistics-card
          title="北京销售商余量"
          :value="Dealer_balance_A"
          :percentage="{
            value: '+3%',
            color: 'text-success',
          }"
          :icon="{
            component: 'ni ni-cart',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div>
      <div class="col-xl-4 col-sm-6 mb-xl-0 mb-4">
        <div style="height: 30px;"></div>
        <mini-statistics-card
          title="上海销售商余量"
          :value="Dealer_balance_B"
          :percentage="{
            value: '-2%',
            color: 'text-danger',
          }"
          :icon="{
            component: 'ni ni-cart',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div>
      <div class="col-xl-4 col-sm-6 mb-xl-0">
        <div style="height: 30px;"></div>
        <mini-statistics-card
          title="天津销售商余量"
          :value="Dealer_balance_C"
          :percentage="{
            value: '+5%',
            color: 'text-success',
          }"
          :icon="{
            component: 'ni ni-cart',
            background: iconBackground,
          }"
          direction-reverse
        />
      </div>
    </div>
    <div style="height: 80px;"></div>
    <div class="row">
      <div class="col-lg-7 mb-lg-0 mb-4">
        <div class="card">
          <div class="card-body p-3">
            <div class="row">
              <div class="col-lg-6">
                <div style="height: 0px;"></div>
                <div class="d-flex flex-column h-100">
                  <profile-info-card
                    title="个人资料"
                    
                    :info="{
                      user_name: '一般用户',
                      user_balance:User_balance,
                      wallet_address: wallet_address,
                      user_pub_key: user_pub_key,
                      user_identity: user_identity,
                    }"
                  />
                </div>
              </div>
              <div class="col-lg-5 ms-auto text-center mt-5 mt-lg-0">
                <div class="border-radius-lg h-100">
                  <div
                    class="position-relative d-flex align-items-center justify-content-center h-100"
                  >
                    <img
                      class="w-100 position-relative z-index-2 pt-4"
                      src="../assets/img/illustrations/maotai_.png"
                      alt="rocket"
                      style="border-radius: 10px;"                    />
                  </div>
                </div>
              </div>
              <div class="col-lg-5 ms-auto text-center mt-5 mt-lg-0">

              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="col-lg-5">
        <div class="card h-100 p-3" style="position: relative;">

          <div class="col-lg-6 col-md-3">
            <timeline-list
              class="h-100"
              title="当前货品溯源情况"
              description="通过右侧按钮获取最新信息"
            >
            <TimelineItem
              v-for="(item, index) in traceInfo"
              :key="index"
              color="success"
              icon="check-bold"
              :title="item.Description"
              :date-time="item.Time"
            />
            </timeline-list>
            <div class="right-side-content">

              <div class="child-div" style="padding-top: 100px;padding-left: 50px;">
                <soft-button color="dark" variant="gradient" @click="UserBuy">
                  <i class="fas fa-plus me-2"></i>
                  购买茅台
                </soft-button>
              </div>
              <div class="child-div" style="padding-top: 50px;padding-left: 45px;">
                <soft-button color="dark" variant="gradient" @click="ProducerBuy">
                  <i class="fas fa-plus me-2"></i>
                  生产商进货
                </soft-button>
              </div>
              <div class="child-div" style="padding-bottom:50px;padding-left: 45px;">
                <soft-button color="dark" variant="gradient" @click="DealerBuy">
                  <i class="fas fa-plus me-2"></i>
                  销售商进货
                </soft-button>
              </div>
          </div>
        
      </div>
        </div>
      </div>
    </div>
    <div class="mt-4 row">
      <div class="mb-4 col-lg-5 mb-lg-0">
        <div class="card z-index-2">
        </div>
      </div>
      <div class="col-lg-7">
        <div class="card z-index-2">

        </div>
      </div>
    </div>
    <div class="row my-4">
      <div class="col-lg-8 col-md-6 mb-md-0 mb-4">
      </div>
      <div class="col-lg-4 col-md-6">
      </div>
    </div>
  </div>
</template>
<script>

// 在这里写后端访问逻辑
import { onMounted, ref } from 'vue'
import axios from 'axios'

// 获取余额的逻辑
const Raw_balance=ref('')
const Producer_balance_A=ref('')
const Producer_balance_B=ref('')
const Dealer_balance_A=ref('')
const Dealer_balance_B=ref('')
const Dealer_balance_C=ref('')
const User_balance=ref('')

const initBalanceData=async()=>{


  axios.get('http://localhost:8081/all_balance')
  .then(response=>{
    Raw_balance.value=response.data.Raw_balance
    Producer_balance_A.value=response.data.A_producer_balance
    Producer_balance_B.value=response.data.B_producer_balance
    Dealer_balance_A.value=response.data.A_dealer_balance
    Dealer_balance_B.value=response.data.B_dealer_balance
    Dealer_balance_C.value=response.data.C_dealer_balance
    User_balance.value=response.data.User_balance
  })
  .catch(err=>{
    console.log(err)
  })
}

// 获取个人信息的逻辑
const wallet_address=ref()
const user_pub_key=ref()
const user_identity=ref()

const getUserInfo=async()=>{


axios.get('http://localhost:8081/wallet_info_ref/%E7%94%A8%E6%88%B7')
.then(response=>{
  wallet_address.value=response.data.Address
  user_pub_key.value=response.data.PublicKey
  user_identity.value=response.data.Identity

})
.catch(err=>{
  console.log(err)
})
}



// 获取溯源信息

const traceInfo=ref()

const getTraceInfo=()=>{
  axios.get('http://localhost:8081/trace_currency')
  .then(response=>{
    traceInfo.value=response.data
  })
  .catch(err=>{
    console.log(err)
  })

}

// import { ElMessageBox } from 'element-plus'

// 购买
const buy_result=ref()

const UserBuy = async () => {
  try {
    const response = await axios.get('http://localhost:8081/buy');
    buy_result.value = response.data;
    
    if (buy_result.value.Success) {
      traceInfo.value=response.data.TraceTrades
      await initBalanceData(); // 等待 initBalanceData 函数执行完毕

    }
  } catch (err) {
    console.log(err);
  }
}

const DealerBuy = async () => {
  try {
    await axios.get('http://localhost:8081/dealer_buy');
    await initBalanceData(); // 等待 initBalanceData 函数执行完毕

  } catch (err) {
    console.log(err);
  }
}
const ProducerBuy = async () => {
  try {
    await axios.get('http://localhost:8081/producer_buy');
    await initBalanceData(); // 等待 initBalanceData 函数执行完毕

  } catch (err) {
    console.log(err);
  }
}


const btn_click=()=>{
  // user_identity.value="sacds"
  // getTraceInfo(); 

}




import ProfileInfoCard from "./components/ProfileInfoCard.vue";
import SoftButton from "@/components/SoftButton.vue";
import MiniStatisticsCard from "@/examples/Cards/MiniStatisticsCard.vue";
// import ReportsBarChart from "@/examples/Charts/ReportsBarChart.vue";
// import GradientLineChart from "@/examples/Charts/GradientLineChart.vue";
import TimelineList from "./components/TimelineList.vue";
import TimelineItem from "./components/TimelineItem.vue";
// import ProjectsCard from "./components/ProjectsCard.vue";
import US from "../assets/img/icons/flags/US.png";
import DE from "../assets/img/icons/flags/DE.png";
import GB from "../assets/img/icons/flags/GB.png";
import BR from "../assets/img/icons/flags/BR.png";
import {
  faHandPointer,
  faUsers,
  faCreditCard,
  faScrewdriverWrench,
} from "@fortawesome/free-solid-svg-icons";
export default {
  name: "dashboard-default",
  setup(){
    onMounted(()=>{

      initBalanceData()
      getUserInfo()
    }

    );
    return{
      btn_click,
      getTraceInfo,
      UserBuy,
      DealerBuy,
      ProducerBuy,
      Raw_balance,
      Producer_balance_A,
      Producer_balance_B,
      Dealer_balance_A,
      Dealer_balance_B,
      Dealer_balance_C,
      User_balance,
      wallet_address,
      user_identity,
      user_pub_key,
      traceInfo
    }
  },
  data() {
    return {
      iconBackground: "bg-gradient-success",
      faCreditCard,
      faScrewdriverWrench,
      faUsers,
      faHandPointer,
      sales: {
        us: {
          country: "United States",
          sales: 2500,
          value: "$230,900",
          bounce: "29.9%",
          flag: US,
        },
        germany: {
          country: "Germany",
          sales: "3.900",
          value: "$440,000",
          bounce: "40.22%",
          flag: DE,
        },
        britain: {
          country: "Great Britain",
          sales: "1.400",
          value: "$190,700",
          bounce: "23.44%",
          flag: GB,
        },
        brasil: {
          country: "Brasil",
          sales: "562",
          value: "$143,960",
          bounce: "32.14%",
          flag: BR,
        },
      },
    };
  },
  components: {
    MiniStatisticsCard,
    // ReportsBarChart,
    // GradientLineChart,
    // ProjectsCard,
    TimelineList,
    TimelineItem,
    ProfileInfoCard,
    SoftButton
  },
};

</script>

<style>
.right-side-content {
  position: absolute;
  right: 0;     /* 贴紧父元素的右侧 */
  top: 0;       /* 从父元素的顶部开始，可以调整 */
  width: 40%;   /* 宽度为父元素的50% */
  height: 100%; /* 如果需要，可以调整高度 */
}

.child-div {
  width: 100%;         /* 子元素宽度等于父元素宽度 */
  height: 33.333%;     /* 子元素高度为父元素高度的三分之一 */
  box-sizing: border-box; /* 如果有边框或内边距，确保它们包含在高度内 */
}
</style>