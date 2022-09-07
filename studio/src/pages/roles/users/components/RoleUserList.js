import React from 'react';
import { Popconfirm, Button, Table } from 'antd';
import { useDispatch } from 'react-redux';
import { DeleteOutlined } from '@ant-design/icons';
import { deleteRoleUser, getRoles } from '../../../../actions/roles';

export default function UserList({ users, roleID }) {
  const dispatch = useDispatch();
  const columns = [
    { title: 'First Name', dataIndex: 'first_name', key: 'name' },
    { title: 'Last Name', dataIndex: 'last_name', key: 'last_name' },
    {
      title: 'Display Name',
      dataIndex: 'display_name',
      key: 'display_name',
      ellipsis: true,
    },
    { title: 'E-mail', dataIndex: 'email', key: 'email' },

    {
      title: 'Action',
      dataIndex: 'operation',
      width: '30%',
      render: (_, record) => {
        return (
          <span>
            <Popconfirm
              title="Sure to Delete?"
              onConfirm={() => {
                deleteUser(record.id);
              }}
            >
              <Button type="danger">
                <DeleteOutlined />
              </Button>
            </Popconfirm>
          </span>
        );
      },
    },
  ];

  const deleteUser = (id) => {
    dispatch(deleteRoleUser(roleID, id)).then(() => dispatch(getRoles()));
  };

  return <Table bordered columns={columns} dataSource={users} rowKey={'id'} />;
}
