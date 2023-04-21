import React from 'react';
import { Card, Avatar, Row, Col } from 'antd';
import { UserOutlined, EyeTwoTone, UsergroupAddOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { Helmet } from 'react-helmet';
const Members = () => {
  const { Meta } = Card;
  const gridStyle = {
    textAlign: 'left',
  };
  return (
    <div>
      <Row gutter={[16, 24]}>
        <Helmet title={'Members'} />
        <Col xs={24} md={12}>
          <Link to="/settings/members/users">
            <Card style={gridStyle} hoverable>
              <Meta
                avatar={
                  <Avatar
                    gap={4}
                    icon={<UserOutlined twoToneColor="#ffb41f" />}
                    style={{ backgroundColor: '#E8EFF2', color: '#ffb41f' }}
                  />
                }
                title="Users"
                description="View User Details"
              />
            </Card>
          </Link>
        </Col>
        <Col xs={24} md={12}>
          <Link to="/settings/members/policies">
            <Card style={gridStyle} hoverable>
              <Meta
                avatar={
                  <Avatar
                    gap={4}
                    icon={<EyeTwoTone twoToneColor="#51bbf6" />}
                    style={{ backgroundColor: '#E8EFF2' }}
                  />
                }
                title="Policies"
                description="Update user policies"
              />
            </Card>
          </Link>
        </Col>
        <Col xs={24} md={12}>
          <Link to="/settings/members/roles">
            <Card style={gridStyle} hoverable>
              <Meta
                avatar={
                  <Avatar
                    gap={4}
                    icon={<UsergroupAddOutlined twoToneColor="#ffb41f" />}
                    style={{ backgroundColor: '#E8EFF2', color: '#ffb41f' }}
                  />
                }
                title="Roles"
                description="Role settings"
              />
            </Card>
          </Link>
        </Col>
      </Row>
    </div>
  );
};

export default Members;
