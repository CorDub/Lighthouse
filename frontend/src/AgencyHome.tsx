import Navbar from "./Navbar.tsx";
import LeftNav from "./LeftNav.tsx";
import './styles/Home.css';

function AgencyHome() {
  return (
    <div class="home">
      <Navbar />
      <LeftNav />
      <p>You don't seem to have any reports yet. </p>
      <p>Create a new report in the vertical menu on the left.</p>
    </div>
  )
}

export default AgencyHome