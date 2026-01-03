import React, { useState } from "react";
import { DepartmentCard } from "../department-card/DepartmentCard";

interface User {
  id: string;
  name: string;
}

interface Department {
  id: string;
  name: string;
  users: User[];
}

interface DepartmentListProps {
  disabled?: boolean;
  title?: string;
  onSelectionChange?: (sel: { departments: string[]; users: string[] }) => void;
}
const DepartmentList: React.FC<DepartmentListProps> = ({
  disabled = false,
  title = "Department List",
  onSelectionChange,
}) => {
  const [search, setSearch] = useState("");
  const departments: Department[] = [
    {
      id: "d1",
      name: "IT Department",
      users: [
        { id: "u1", name: "Khalil" },
        { id: "u2", name: "Adam" },
      ],
    },
    {
      id: "d2",
      name: "HR Department",
      users: [
        { id: "u3", name: "Sara" },
        { id: "u4", name: "Yasmine" },
      ],
    },
  ];

  const [openDept, setOpenDept] = useState<string | null>(null);
  const [checkedDept, setCheckedDept] = useState<string[]>([]);
  const [checkedUsers, setCheckedUsers] = useState<string[]>([]);

  const filteredDepartments = departments.filter((d) =>
    d.name.toLowerCase().includes(search.trim().toLowerCase())
  );

  const toggleDept = (id: string) => {
    setOpenDept(openDept === id ? null : id);
  };

  React.useEffect(() => {
    onSelectionChange && onSelectionChange({ departments: checkedDept, users: checkedUsers });
  }, [checkedDept, checkedUsers, onSelectionChange]);

  return (
    <div
      className={`${disabled ? "opacity-40 pointer-events-none" : ""} 
      w-full max-w-sm p-4 sm:p-6 bg-neutral-primary-soft 
      border-l-indigo-200 rounded shadow-xs`}
    >
      <div className="flex items-center justify-between mb-2">
        <h5 className="text-base md:text-xl font-semibold text-heading">
          {title}
        </h5>
        <div className="ml-3 w-44">
          <div className="relative w-full">
            <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
              <svg className="w-4 h-4 text-gray-500" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 8v8m0-8a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm0 8a2 2 0 1 0 0 4 2 2 0 0 0 0-4Zm8-8a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm0 0a4 4 0 0 1-4 4h-1a3 3 0 0 0-3 3"/>
              </svg>
            </div>
            <input
              type="text"
              id="simple-search"
              className="pl-9 pr-3 py-2.5 bg-gray-950 border border-violet-900 rounded-md text-gray-200 text-sm focus:ring-2 focus:ring-brand focus:border-brand block w-full placeholder-gray-500"
              placeholder="Search branch name"
              required
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>
      </div>

      <p className="text-body">
        Connect with one of our available department list or create a new one.
      </p>


      <ul className="my-6 space-y-3">
        {filteredDepartments.map((dept) => (
          <li key={dept.id}>
            <DepartmentCard
              title={dept.name}
              isOpen={openDept === dept.id}
              checked={checkedDept.includes(dept.id)}
              onToggle={() => toggleDept(dept.id)}
              onCheck={() =>
                setCheckedDept((prev) =>
                  prev.includes(dept.id)
                    ? prev.filter((d) => d !== dept.id)
                    : [...prev, dept.id]
                )
              }
            />

            {openDept === dept.id && (
              <div className="ml-14 mt-2 space-y-2">
                {dept.users.map((user) => (
                  <label
                    key={user.id}
                    className="flex items-center gap-2 text-sm text-white opacity-80"
                  >
                    <input
                      type="checkbox"
                      checked={checkedUsers.includes(user.id)}
                      onChange={() =>
                        setCheckedUsers((prev) =>
                          prev.includes(user.id)
                            ? prev.filter((u) => u !== user.id)
                            : [...prev, user.id]
                        )
                      }
                    />
                    {user.name}
                  </label>
                ))}
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default DepartmentList;
