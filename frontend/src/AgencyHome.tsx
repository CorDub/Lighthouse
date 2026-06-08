import Navbar from "./Navbar.tsx";
import LeftNav from "./LeftNav.tsx";
import './styles/Home.css';
import OpenPlan from "./OpenPlan.tsx";
import { displayed, type AgencyHomeModule } from "./stores/agencyHomeDisplay.tsx";
import EmptyAgency from "./EmptyAgency.tsx";
import CreateReportModule from "./CreateReportModule.tsx";
import { type Component } from "solid-js";

function AgencyHome() {
  const modules: Record<AgencyHomeModule, Component> = {
    empty: EmptyAgency,
    createReport: CreateReportModule,
  }

  return (
    <div class="home">
      <Navbar />
      <LeftNav />
      <OpenPlan comp={modules[displayed()]}/>
    </div>
  )
}

export default AgencyHome