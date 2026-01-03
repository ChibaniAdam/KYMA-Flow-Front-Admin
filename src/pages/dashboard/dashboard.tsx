import  { useState } from "react";
import RepoList from "../../components/repo-list/repo-list";
import DepartmentList from "../../components/department-list/DepartmentList";

const Dashboard = () => {
    const [firstSelected, setFirstSelected] = useState<string | null>(null);

    return (
        <div>
            <h1>Dashboard</h1>

           <ol className="flex items-center w-full px-4 lg:px-8">
  <li className="flex-1 flex items-center text-fg-brand after:content-[''] after:flex-1 after:h-1 after:border-b after:border-brand-subtle after:border-4 after:ms-4 after:rounded-full">
    <span className="flex items-center justify-center w-10 h-10 bg-brand-softer rounded-full lg:h-12 lg:w-12 shrink-0">
      1 
    </span>
  </li>

  <li className=" flex items-center justify-end">
    <span
      className={`flex items-center justify-center w-10 h-10 ${
        firstSelected ? "bg-brand-softer" : "bg-ne  utral-tertiary"
      } rounded-full lg:h-12 lg:w-12 shrink-0`}
    >
      2
    </span>
  </li>
</ol>


     <div className="w-full mt-6">
  <div className="grid grid-cols-3 gap-4 w-full">
    
    <div>
      <RepoList onSelect={(id) => setFirstSelected(id)} />
    </div>

    <div>
      <DepartmentList disabled={!firstSelected} />
    </div>

    <div>
      <DepartmentList disabled={!firstSelected} />
    </div>

  </div>
</div>
        </div>

        
    )

}

export default Dashboard;