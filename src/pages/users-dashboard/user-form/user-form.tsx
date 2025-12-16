import React, { useEffect, useRef } from "react";
import type { CreateUserInput, UpdateUserInput, User } from "../../../GQL/models/user";
import type { Department } from "../../../GQL/models/department";
import "./user-form.css";
import { CustomSelect } from "../../../components/custom-select/custom-select";

interface Props {
  editingUser: User | null;
  formData: CreateUserInput | UpdateUserInput;
  setFormData: (d: any) => void;
  onSubmit: () => void;
  onClose: () => void;
  departments: Department[] | [];
}

export const UserForm = ({
  editingUser,
  formData,
  setFormData,
  onSubmit,
  onClose,
  departments
}: Props) => {


  const modalRef = useRef<HTMLDivElement>(null);



  useEffect(() => {
    if (!editingUser) {
      const first = formData.givenName?.trim().toLowerCase();
      const last = formData.sn?.trim().toLowerCase();
      if (first && last) {
        setFormData({ ...formData, cn: `${first}.${last}` });
      }
    }
  }, [formData.givenName, formData.sn]);

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    globalThis.addEventListener("keydown", handleEsc);
    return () => globalThis.removeEventListener("keydown", handleEsc);
  }, []);

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget) onClose();
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleDepartmentChange = (value: string) => {
    setFormData({...formData, department: value})
  }

  return (
    <div className="modal-backdrop" onClick={handleBackdropClick}>
      <div className="modal-card modal-3parts" ref={modalRef}>

        <div className="modal-header">
          <h2>{editingUser ? "Update User" : "Create User"}</h2>
          <p>{editingUser ? "Modify the user details below." : "Fill in the information below to add a new user."}</p>
        </div>

        <div className="user-modal-body">

          <div className="form-group">
            <label htmlFor="uid">UID</label>
            <input
              id="uid"
              name="uid"
              value={formData.uid || ""}
              onChange={handleChange}
              placeholder="Enter uid"
            />
          </div>

          <div className="form-group">
            <label htmlFor="givenName">First Name</label>
            <input
              id="givenName"
              name="givenName"
              value={formData.givenName || ""}
              onChange={handleChange}
              placeholder="Enter first name"
            />
          </div>

          <div className="form-group">
            <label htmlFor="sn">Last Name</label>
            <input
              id="sn"
              name="sn"
              value={formData.sn || ""}
              onChange={handleChange}
              placeholder="Enter last name"
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              name="mail"
              value={formData.mail || ""}
              onChange={handleChange}
              placeholder="email@example.com"
            />
          </div>

          <div className="form-group">
            <label htmlFor="department">Department</label>
                         <CustomSelect
                  value={formData.department || ""}
                  options={departments.map((d) => ({ label: d.ou, value: d.ou }))}
                  placeholder="All Departments"
                  onChange={(v) => handleDepartmentChange(v)}
                />
          </div>

        </div>

        <div className="modal-footer">
          <button className="cancel-btn" onClick={onClose}>
            Cancel
          </button>

          <button className="submit-btn" onClick={onSubmit}>
            {editingUser ? "Update" : "Create"}
          </button>
        </div>

      </div>
    </div>
  );
};
