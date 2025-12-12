import React, { useEffect, useState } from "react";
import {
  listUsers,
  createUser,
  updateUser,
  deleteUser,
} from "../../services/userService";
import type { User, CreateUserInput, UpdateUserInput } from "../../GQL/models/user";
import "./users-dashboard.css";

export const UsersDashboard = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
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

  // Load users
  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await listUsers();
      setUsers(data);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

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
    setFormData({
      uid: user.uid,
      cn: user.cn,
      sn: user.sn,
      givenName: user.givenName,
      mail: user.mail,
      department: user.department,
      password: "",
      repositories: user.repositories,
    });
    setShowModal(true);
  };

  const handleSubmit = async () => {
    try {
      if (editingUser) {
        await updateUser(formData as UpdateUserInput);
      } else {
        await createUser(formData as CreateUserInput);
      }
      setShowModal(false);
      fetchUsers();
    } catch (err) {
      console.error(err);
    }
  };

  const handleDelete = async (uid: string) => {
    if (window.confirm("Are you sure you want to delete this user?")) {
      await deleteUser(uid);
      fetchUsers();
    }
  };

  return (
    <div className="user-table-container">
      <div className="table-header">
        <h2>Users</h2>
        <button className="create-btn" onClick={handleCreateClick}>
          Create User
        </button>
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <table className="user-table">
          <thead>
            <tr>
              <th>UID</th>
              <th>Name</th>
              <th>Email</th>
              <th>Department</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.uid}>
                <td>{user.uid}</td>
                <td>{user.givenName} {user.sn}</td>
                <td>{user.mail}</td>
                <td>{user.department}</td>
                <td>
                  <button
                    className="update-btn"
                    onClick={() => handleEditClick(user)}
                  >
                    Update
                  </button>
                  <button
                    className="delete-btn"
                    onClick={() => handleDelete(user.uid)}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {showModal && (
        <div className="modal-backdrop">
          <div className="modal">
            <h3>{editingUser ? "Update User" : "Create User"}</h3>
            <input
              name="uid"
              placeholder="UID"
              value={formData.uid}
              onChange={handleChange}
              disabled={!!editingUser}
            />
            <input
              name="givenName"
              placeholder="First Name"
              value={formData.givenName}
              onChange={handleChange}
            />
            <input
              name="sn"
              placeholder="Last Name"
              value={formData.sn}
              onChange={handleChange}
            />
            <input
              name="cn"
              placeholder="Common Name"
              value={formData.cn}
              onChange={handleChange}
            />
            <input
              name="mail"
              placeholder="Email"
              value={formData.mail}
              onChange={handleChange}
            />
            <input
              name="department"
              placeholder="Department"
              value={formData.department}
              onChange={handleChange}
            />
            <input
              name="password"
              type="password"
              placeholder="Password"
              value={formData.password}
              onChange={handleChange}
            />
            <div className="modal-actions">
              <button onClick={handleSubmit}>
                {editingUser ? "Update" : "Create"}
              </button>
              <button onClick={() => setShowModal(false)}>Cancel</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
