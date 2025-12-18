import { useEffect, useState } from "react";
import {
  listUsers,
  createUser,
  updateUser,
  deleteUser,
} from "../../services/userService";

import type { User, CreateUserInput, UpdateUserInput, UserPage } from "../../GQL/models/user";
import { UserForm } from "./user-form/user-form";
import { listDepartments } from "../../services/departmentService";
import type { Department } from "../../GQL/models/department";

import "./users-dashboard.css";
import { ConfirmationModal } from "../../components/confirmation-modal/confirmation-modal";
import { FilterBar } from "../../components/filter-bar/filter-bar";
import { DataTable } from "../../components/data-table/data-table";

export const UsersDashboard = () => {
  const [usersPage, setUsersPage] = useState<UserPage | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(8);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [selected, setSelected] = useState<string[]>([]);

  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const [deleteUserId, setDeleteUserId] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  const [search, setSearch] = useState("");
  const [searchByMail, setSearchByMail] = useState("");

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
    try {
      const filter: any = {};
      if (departmentFilter) {filter.department = departmentFilter;}
      if (search.trim()) filter.cn = search.trim().replace(" ",".");
      if (searchByMail.trim()) filter.mail = searchByMail.trim();

      const data = await listUsers(filter, {
        page,
        limit: pageSize,
      });

      setUsersPage(data);
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
  fetchUsers();
  }, [page, search, departmentFilter, searchByMail]);



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

        <FilterBar
            filters={[
              {
                key: "search",
                type: "text",
                placeholder: "Search by name",
                value: search,
                onChange: setSearch,
              },
              {
                key: "searchMail",
                type: "text",
                placeholder: "Search by email...",
                value: searchByMail,
                onChange: setSearchByMail,
              },
              {
                key: "department",
                type: "select",
                placeholder: "All Departments",
                value: departmentFilter,
                options: departments.map((d) => ({
                  label: d.ou,
                  value: d.ou,
                })),
                onChange: setDepartmentFilter,
              },
            ]}
            actions={
              <button className="create-btn" onClick={handleCreateClick}>
                + Add User
              </button>
            }
          />

 <div className="table-wrapper">
<DataTable<User>
  data={usersPage?.items ?? []}
  loading={loading}
  page={page}
  pageSize={pageSize}
  total={usersPage?.total ?? 0}
  onPageChange={setPage}
  columns={[
    { key: "uid", header: "UID", sortable: true },
    {
      key: "givenName",
      header: "Name",
      sortable: true,
      render: (u) => `${u.givenName} ${u.sn}`,
    },
    { key: "mail", header: "Email", sortable: true },
    { key: "department", header: "Department", sortable: true },
  ]}
  onEdit={handleEditClick}
  onDelete={(u) => handleDelete(u.uid)}
/>

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
