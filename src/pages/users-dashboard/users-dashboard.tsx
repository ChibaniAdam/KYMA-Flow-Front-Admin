import { useEffect, useState } from "react";
import {
  listUsers,
  createUser,
  updateUser,
  deleteUser,
} from "../../services/userService";

import type { User, CreateUserInput, UpdateUserInput } from "../../GQL/models/user";
import { UserForm } from "./user-form/user-form";
import { listDepartments } from "../../services/departmentService";
import type { Department } from "../../GQL/models/department";

import "./users-dashboard.css";
import { CustomSelect } from "../../components/custom-select/custom-select";
import { ConfirmationModal } from "../../components/confirmation-modal/confirmation-modal";

export const UsersDashboard = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);

  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const [deleteUserId, setDeleteUserId] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  const [search, setSearch] = useState("");
  const [departmentFilter, setDepartmentFilter] = useState("");

  const [formData, setFormData] = useState<CreateUserInput | UpdateUserInput>({
    uid: "",
    cn: "",
    sn: "",
    givenName: "",
    mail: "",
    department: "",
    password: "",
    repositories: [],
  });

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await listUsers();
      setUsers(data);
      setFilteredUsers(data);
    } finally {
      setLoading(false);
    }
  };

  const fetchDepartments = async () => {
    const data = await listDepartments();
    setDepartments(data);
  };

  useEffect(() => {
    fetchUsers();
    fetchDepartments();
  }, []);


  useEffect(() => {
    let u = [...users];

    if (departmentFilter)
      u = u.filter((usr) => usr.department === departmentFilter);

    if (search.trim()) {
      const s = search.toLowerCase();
      u = u.filter(
        (usr) =>
          usr.uid.toLowerCase().includes(s) ||
          usr.givenName.toLowerCase().includes(s) ||
          usr.sn.toLowerCase().includes(s) ||
          usr.mail.toLowerCase().includes(s)
      );
    }

    setFilteredUsers(u);
  }, [search, departmentFilter, users]);


  const handleCreateClick = () => {
    setEditingUser(null);
    setFormData({
      uid: "",
      cn: "",
      sn: "",
      givenName: "",
      mail: "",
      department: "",
      password: "",
      repositories: [],
    });
    setShowModal(true);
  };

  const handleEditClick = (user: User) => {
    setEditingUser(user);
    setFormData(user);
    setShowModal(true);
  };

  const handleSubmit = async () => {
    try {
      if (editingUser)
        await updateUser(formData as UpdateUserInput);
      else
        await createUser(formData as CreateUserInput);

      setShowModal(false);
      fetchUsers();
    } catch (err) {
      console.error(err);
    }
  };



  const handleDelete = (uid: string) => {
    setDeleteUserId(uid); 
  };

  const confirmDelete = async () => {
    if (!deleteUserId) return;
    setDeleting(true);
    try {
      await deleteUser(deleteUserId);
      fetchUsers();
    } finally {
      setDeleting(false);
      setDeleteUserId(null);
    }
  };

  return (
    <div className="users-page">

      <div className="users-page-title">
        <h1>Users</h1>
        <p>Manage all user accounts across departments.</p>
      </div>

      <div className="users-filter-bar">
        <div className="filter-group">
          <input
            type="text"
            placeholder="Search by name, UID, email..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="filter-input"
          />
        </div>
        <div className="flex gap-2">
        <div className="filter-group select-wrapper">
              <CustomSelect
      value={departmentFilter}
      options={departments.map((d) => ({ label: d.ou, value: d.ou }))}
      placeholder="All Departments"
      onChange={(v) => setDepartmentFilter(v)}
    />
 
        </div>

        <button className="create-btn" onClick={handleCreateClick}>
          + Add User
        </button>
        </div>
      </div>

      <div className="user-table-wrapper">
        {loading ? (
          <div className="skeleton-table">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="skeleton-row"></div>
            ))}
          </div>
        ) : (
          <table className="user-table">
            <thead>
              <tr>
                <th>UID</th>
                <th>Name</th>
                <th>Email</th>
                <th>Department</th>
                <th></th>
              </tr>
            </thead>

            <tbody>
              {filteredUsers.map((user) => (
                <tr key={user.uid}>
                  <td>{user.uid}</td>
                  <td>{user.givenName} {user.sn}</td>
                  <td>{user.mail}</td>
                  <td>{user.department}</td>
                  <td className="actions">
                    <button
                      onClick={() => handleEditClick(user)}
                      className="update-btn"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(user.uid)}
                      className="delete-btn"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}

              {filteredUsers.length === 0 && (
                <tr>
                  <td colSpan={5} className="empty-msg">
                    No users found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>

      {showModal && (
        <UserForm
          editingUser={editingUser}
          formData={formData}
          setFormData={setFormData}
          onSubmit={handleSubmit}
          departments={departments}
          onClose={() => setShowModal(false)}
        />
      )}
      {deleteUserId && (
        <ConfirmationModal
          message="Are you sure you want to delete this user? This action cannot be undone."
          onCancel={() => setDeleteUserId(null)}
          onConfirm={confirmDelete}
          loading={deleting}
        />
      )}

    </div>
  );
};
