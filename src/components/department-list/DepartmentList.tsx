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

  const toggleDept = (id: string) => {
    setOpenDept(openDept === id ? null : id);
  };

  // Notify parent when selection changes
  React.useEffect(() => {
    onSelectionChange && onSelectionChange({ departments: checkedDept, users: checkedUsers });
  }, [checkedDept, checkedUsers, onSelectionChange]);

  return (
    <div
      className={`${disabled ? "opacity-40 pointer-events-none" : ""} 
      w-full max-w-sm p-4 sm:p-6 bg-neutral-primary-soft 
      border-l-indigo-200 rounded shadow-xs`}
    >
     
      <h5 className="mb-2 text-base md:text-xl font-semibold text-heading">
        {title}
      </h5>

      <p className="text-body">
        Connect with one of our available department list or create a new one.
      </p>


      <ul className="my-6 space-y-3">
        {departments.map((dept) => (
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

            {/* USERS DROPDOWN */}
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
