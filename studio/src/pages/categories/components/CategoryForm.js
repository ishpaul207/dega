import React, { useState } from 'react';
import { Row, Col, Button, Form, Input, Space, Switch, Collapse } from 'antd';
import { maker, checker } from '../../../utils/sluger';
import MediaSelector from '../../../components/MediaSelector';
import Editor from '../../../components/Editor';
import Selector from '../../../components/Selector';
import MonacoEditor from '../../../components/MonacoEditor';
import getJsonValue from '../../../utils/getJsonValue';
import SocialCardPreview from '../../../components/PreviewSocialCard';
import PlaceholderImage from '../../../components/ErrorsAndImage/PlaceholderImage';

const CategoryForm = ({ onCreate, data = {} }) => {
  const meta = {
    title: 'CATEGORY NAME',
    canonical_URL: '/category',
    description:
      'Lorem ipsum dolor sit amet consectetur, adipisicing elit. Alias nisi inventore delectus amet tempora voluptatibus sed quibusdam aperiam fugit explicabo consectetur blanditiis cum ab fugiat perspiciatis, eum minima rerum ipsa?',
    // image: PlaceholderImage,
  };
  if (data && data.meta_fields) {
    if (typeof data.meta_fields !== 'string') {
      data.meta_fields = JSON.stringify(data.meta_fields);
    }
  }
  const [form] = Form.useForm();
  const [valueChange, setValueChange] = React.useState(false);

  const onReset = () => {
    form.resetFields();
  };

  const onTitleChange = (string) => {
    form.setFieldsValue({
      slug: maker(string),
    });
  };

  return (
    <Form
      form={form}
      initialValues={{ ...data }}
      name="create-category"
      layout="vertical"
      onFinish={(values) => {
        if (values.meta_fields) {
          values.meta_fields = getJsonValue(values.meta_fields);
        }
        onCreate(values);
        onReset();
      }}
      onValuesChange={() => {
        setValueChange(true);
      }}
    >
      <Row justify="center" style={{ maxWidth: '1200px', width: '100%', margin: '0 auto' }}>
        <Col span={24}>
          <Row justify="end" gutter={40}>
            <Form.Item>
              <Space>
                <Button htmlType="button" onClick={onReset}>
                  Reset
                </Button>
                <Button disabled={!valueChange} type="primary" htmlType="submit">
                  {data && data.id ? 'Update' : 'Submit'}
                </Button>
              </Space>
            </Form.Item>
          </Row>
        </Col>
        <Col span={24}>
          <Row
            gutter={40}
            justify="space-around"
            style={{ background: '#f0f2f5', padding: '1.25rem', marginBottom: '1rem' }}
          >
            <Col span={12}>
              <Form.Item
                name="name"
                label="Category Name"
                rules={[
                  {
                    required: true,
                    message: 'Please enter the name!',
                  },
                  { min: 3, message: 'Name must be minimum 3 characters.' },
                  { max: 50, message: 'Name must be maximum 50 characters.' },
                ]}
              >
                <Input onChange={(e) => onTitleChange(e.target.value)} />
              </Form.Item>
              <Row gutter={40}>
                <Col md={{ span: 16 }}>
                  <Form.Item name="parent_id" label="Parent Category">
                    <Selector action="Categories" />
                  </Form.Item>
                </Col>
                <Col md={{ span: 5 }}>
                  <Form.Item label="Featured" name="is_featured" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                name="slug"
                label="Slug"
                rules={[
                  {
                    required: true,
                    message: 'Please input the slug!',
                  },
                  {
                    pattern: checker,
                    message: 'Please enter valid slug!',
                  },
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
            <Col span={10}>
              <Form.Item label="Featured Image" name="medium_id">
                <MediaSelector maxWidth={'350px'} containerStyles={{ justifyContent: 'start' }} />
              </Form.Item>
            </Col>
            <Col span={12} style={{ marginRight: 'auto', marginLeft: '20px' }}>
              <Form.Item name="description" label="Description">
                <Editor
                  style={{ width: '600px', background: '#fff', padding: '0.5rem 0.75rem' }}
                  placeholder="Enter Description..."
                  basic={true}
                />
              </Form.Item>
            </Col>
          </Row>
        </Col>
        <Col span={24}>
          <Row gutter={40} style={{ background: '#f0f2f5' }}>
            <Collapse
              expandIconPosition="right"
              expandIcon={({ isActive }) => <Button>{isActive ? 'Close' : 'Expand'}</Button>}
              style={{ width: '100%' }}
            >
              <Collapse.Panel header="Meta Data" className="meta-data-container">
                <div style={{ flexBasis: '67%', marginRight: '2.4rem' }}>
                  <Form.Item name={['meta', 'title']} label="Meta Title">
                    <Input />
                  </Form.Item>
                  <Form.Item name={['meta', 'description']} label="Meta Description">
                    <Input.TextArea />
                  </Form.Item>
                  <Form.Item name={['meta', 'canonical_URL']} label="Canonical URL">
                    <Input />
                  </Form.Item>
                </div>
                <SocialCardPreview
                  type="google"
                  url={form.getFieldValue(['meta', 'canonical_URL'])}
                  title={meta.title}
                  desc={meta.description}
                />
              </Collapse.Panel>
              <Collapse.Panel header="Twitter Card" className="meta-data-container">
                <div style={{ flexBasis: '67%', marginRight: '2.4rem' }}>
                  <Form.Item label="Twitter Image" name={['meta', 'medium_id']}>
                    <MediaSelector
                      maxWidth={'350px'}
                      containerStyles={{ justifyContent: 'start' }}
                    />
                  </Form.Item>
                  <Form.Item name={['meta', 'title']} label="Meta Title">
                    <Input />
                  </Form.Item>
                  <Form.Item name={['meta', 'description']} label="Meta Description">
                    <Input.TextArea />
                  </Form.Item>
                  <Form.Item name={['meta', 'canonical_URL']} label="Canonical URL">
                    <Input />
                  </Form.Item>
                </div>
                <SocialCardPreview
                  type="twitter"
                  url={meta.canonical_URL}
                  title={meta.title}
                  desc={meta.description}
                  image={meta.image}
                />
              </Collapse.Panel>
              <Collapse.Panel header="Facebook Card" className="meta-data-container">
                <div style={{ flexBasis: '67%', marginRight: '2.4rem' }}>
                  <Form.Item label="Facebook Image" name={['meta', 'medium_id']}>
                    <MediaSelector
                      maxWidth={'350px'}
                      containerStyles={{ justifyContent: 'start' }}
                    />
                  </Form.Item>
                  <Form.Item name={['meta', 'title']} label="Meta Title">
                    <Input />
                  </Form.Item>
                  <Form.Item name={['meta', 'description']} label="Meta Description">
                    <Input.TextArea />
                  </Form.Item>
                  <Form.Item name={['meta', 'canonical_URL']} label="Canonical URL">
                    <Input />
                  </Form.Item>
                </div>
                <SocialCardPreview
                  type="fb"
                  url={meta.canonical_URL}
                  title={meta.title}
                  desc={meta.description}
                  image={meta.image}
                />
              </Collapse.Panel>
            </Collapse>
            <Collapse
              expandIconPosition="right"
              expandIcon={({ isActive }) => <Button>{isActive ? 'Close' : 'Expand'}</Button>}
              style={{ width: '100%' }}
            >
              <Collapse.Panel header="Code Injection">
                <Form.Item name="header_code" label="Header Code">
                  <MonacoEditor language="html" width="100%" />
                </Form.Item>
                <Form.Item name="footer_code" label="Footer Code">
                  <MonacoEditor language="html" width="100%" />
                </Form.Item>
              </Collapse.Panel>
            </Collapse>
            <Collapse
              expandIconPosition="right"
              expandIcon={({ isActive }) => <Button>{isActive ? 'Close' : 'Expand'}</Button>}
              style={{ width: '100%' }}
            >
              <Collapse.Panel header="Meta Fields">
                <Form.Item name="meta_fields">
                  <MonacoEditor language="json" width="100%" />
                </Form.Item>
              </Collapse.Panel>
            </Collapse>
          </Row>
        </Col>
      </Row>
    </Form>
  );
};

export default CategoryForm;
