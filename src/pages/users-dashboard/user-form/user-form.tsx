import { useEffect } from "react";
import type { CreateUserInput, UpdateUserInput, User } from "../../../GQL/models/user";
import type { Department } from "../../../GQL/models/department";
import { CustomSelect } from "../../../components/custom-select/custom-select";
import "./user-form.css";

interface Props {
  editingUser: User | null;
  formData: CreateUserInput | UpdateUserInput;
  setFormData: (d: any) => void;
  departments: Department[];
}

export const UserForm = ({
  editingUser,
  formData,
  setFormData,
  departments,
}: Props) => {
  useEffect(() => {
    if (!editingUser) {
      const first = formData.givenName?.trim().toLowerCase();
      const last = formData.sn?.trim().toLowerCase();
      if (first && last) {
        setFormData({ ...formData, cn: `${first}.${last}` });
      }
    }
  }, [formData.givenName, formData.sn]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  return (
    <>
      <div className="form-group">
        <label htmlFor="uid">UID</label>
        <input name="uid" value={formData.uid || ""} onChange={handleChange} />
      </div>

      <div className="form-group">
        <label htmlFor="givenName">First Name</label>
        <input
          name="givenName"
          value={formData.givenName || ""}
          onChange={handleChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="sn">Last Name</label>
        <input
          name="sn"
          value={formData.sn || ""}
          onChange={handleChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="mail">Email</label>
        <input
          type="email"
          name="mail"
          value={formData.mail || ""}
          onChange={handleChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="department">Department</label>
        <CustomSelect
          value={formData.department || ""}
          options={departments.map(d => ({
            label: d.ou,
            value: d.ou,
          }))}
          onChange={(v) =>
            setFormData({ ...formData, department: v })
          }
        />
      </div>
    </>
  );
};
